package fscli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/c-bata/go-prompt"
	"github.com/olekukonko/tablewriter"
	"github.com/shibukawa/configdir"
	"golang.org/x/exp/slices"
)

const LongLine = "--------------------------------------------------------------------------"

type OutputMode string

const (
	OutputModeJSON  OutputMode = "json"
	OutputModeTable OutputMode = "table"
)

const (
	VENDOR_NAME  = "maruware"
	APP_NAME     = "fscli"
	HISTORY_FILE = "history"
)

type Repl struct {
	ctx              context.Context
	fs               *firestore.Client
	in               io.Reader
	out              io.Writer
	outputMode       OutputMode
	exe              *Executor
	enabledPager     bool
	collectionsCache map[string][]string
}

func NewRepl(ctx context.Context, fs *firestore.Client, in io.Reader, out io.Writer, outputMode OutputMode) *Repl {
	return &Repl{
		ctx:              ctx,
		fs:               fs,
		in:               in,
		out:              out,
		outputMode:       outputMode,
		exe:              NewExecutor(ctx, fs),
		enabledPager:     false,
		collectionsCache: map[string][]string{},
	}
}

func (r *Repl) completer(d prompt.Document) []prompt.Suggest {
	w := d.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}

	text := d.TextBeforeCursor()

	findCollections := func(baseDoc string) ([]string, error) {
		fn := func(baseDoc string) (*firestore.CollectionIterator, error) {
			return getCollectionsIterator(r.ctx, r.fs, baseDoc)
		}
		collections := getCollections(baseDoc, fn)
		return collections, nil
	}

	c := NewCompleter(NewLexer(text), findCollections)
	suggestions, err := c.Parse()
	if err != nil {
		return []prompt.Suggest{}
	}

	return suggestions
}

func (r *Repl) Start() {
	history := r.readHistory()

	p := prompt.New(
		r.promptProcessLine,
		r.completer,
		prompt.OptionPrefix("> "),
		prompt.OptionSwitchKeyBindMode(prompt.CommonKeyBind),
		prompt.OptionAddASCIICodeBind(prompt.ASCIICodeBind{
			ASCIICode: []byte{0x1b, 0x62}, // Alt/Option + Left
			Fn:        prompt.GoLeftWord,
		}),
		prompt.OptionAddASCIICodeBind(prompt.ASCIICodeBind{
			ASCIICode: []byte{0x1b, 0x66}, // Alt/Option + Right
			Fn:        prompt.GoRightWord,
		}),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlW,
			Fn:  prompt.DeleteWord,
		}),
		prompt.OptionHistory(history),
	)
	p.Run()
}

func (r *Repl) promptProcessLine(line string) {
	if strings.TrimSpace(line) == "" {
		return
	}

	err := r.writeHistory(line)
	if err != nil {
		fmt.Fprintf(r.out, "error: %s\n", err)
		return
	}
	r.ProcessLine(line)
}

func (r *Repl) ProcessLine(line string) {
	lexer := NewLexer(line)
	parser := NewParser(lexer)
	op, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(r.out, "error: %s\n", err)
		return
	}
	if op == nil {
		return
	}

	if err := r.executeOperation(op); err != nil {
		fmt.Fprintf(r.out, "error: %s\n", err)
	}
}

func (r *Repl) executeOperation(op any) error {
	switch v := op.(type) {
	case *MetacommandPager:
		return r.handlePager(v)
	case *MetacommandListCollections:
		return r.handleListCollections(v)
	case *QueryOperation:
		return r.handleQuery(v)
	case *GetOperation:
		return r.handleGet(v)
	case *CountOperation:
		return r.handleCount(v)
	default:
		return fmt.Errorf("unknown operation type")
	}
}

func (r *Repl) handlePager(op *MetacommandPager) error {
	r.enabledPager = op.on
	return nil
}

func (r *Repl) handleListCollections(op *MetacommandListCollections) error {
	cols, err := r.exe.ExecuteListCollections(r.ctx, op)
	if err != nil {
		return err
	}

	out, render := r.pagerableOut()
	for _, col := range cols {
		fmt.Fprintf(out, "%s\n", col)
	}
	return render()
}

func (r *Repl) handleQuery(op *QueryOperation) error {
	docs, err := r.exe.ExecuteQuery(r.ctx, op)
	if err != nil {
		return err
	}

	if r.outputMode == OutputModeJSON {
		r.outputDocsJSON(docs)
	} else if r.outputMode == OutputModeTable {
		r.outputDocsTable(docs)
	}
	return nil
}

func (r *Repl) handleGet(op *GetOperation) error {
	doc, err := r.exe.ExecuteGet(r.ctx, op)
	if err != nil {
		return err
	}

	if r.outputMode == OutputModeJSON {
		r.outputDocJSON(doc)
	} else if r.outputMode == OutputModeTable {
		r.outputDocTable(doc)
	}
	return nil
}

func (r *Repl) handleCount(op *CountOperation) error {
	count, err := r.exe.ExecuteCount(r.ctx, op)
	if err != nil {
		return err
	}

	fmt.Fprintf(r.out, "%d\n", count)
	return nil
}

func (r *Repl) ProcessLineFromPipe() {
	scanner := bufio.NewScanner(r.in)
	for scanner.Scan() {
		line := scanner.Text()
		r.ProcessLine(line)
	}
}


func (r *Repl) outputDocsTable(docs []*firestore.DocumentSnapshot) {
	out, render := r.pagerableOut()
	defer render()

	// collect keys
	keys := []string{}
	for _, doc := range docs {
		for k := range doc.Data() {
			if !slices.Contains(keys, k) {
				keys = append(keys, k)
			}
		}
	}
	slices.Sort(keys)

	table := tablewriter.NewWriter(out)
	table.SetAutoFormatHeaders(false)
	table.SetHeader(append([]string{"ID"}, keys...))

	for _, doc := range docs {
		row := []string{doc.Ref.ID}
		for _, k := range keys {
			val, ok := doc.Data()[k]
			row = append(row, r.toTableCell(val, ok))
		}
		table.Append(row)
	}
	table.Render()
}

func (r *Repl) outputDocTable(doc *firestore.DocumentSnapshot) {
	keys := []string{}
	for k := range doc.Data() {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	table := tablewriter.NewWriter(r.out)
	table.SetAutoFormatHeaders(false)
	table.SetHeader(append([]string{"ID"}, keys...))

	row := []string{doc.Ref.ID}
	for _, k := range keys {
		val, ok := doc.Data()[k]
		row = append(row, r.toTableCell(val, ok))
	}
	table.Append(row)
	table.Render()
}

func (r *Repl) outputDocsJSON(docs []*firestore.DocumentSnapshot) {
	type docOutput struct {
		ID   string         `json:"id"`
		Data map[string]any `json:"data"`
	}

	outputs := make([]docOutput, 0, len(docs))
	for _, doc := range docs {
		outputs = append(outputs, docOutput{
			ID:   doc.Ref.ID,
			Data: doc.Data(),
		})
	}

	j, err := json.Marshal(outputs)
	if err != nil {
		fmt.Fprintf(r.out, "invalid data: %s\n", err)
		return
	}
	fmt.Fprintln(r.out, string(j))
}

func (r *Repl) outputDocJSON(doc *firestore.DocumentSnapshot) {
	output := struct {
		ID   string         `json:"id"`
		Data map[string]any `json:"data"`
	}{
		ID:   doc.Ref.ID,
		Data: doc.Data(),
	}

	j, err := json.Marshal(output)
	if err != nil {
		fmt.Fprintf(r.out, "invalid data: %s\n", err)
		return
	}
	fmt.Fprintln(r.out, string(j))
}

func (r *Repl) toTableCell(val any, ok bool) string {
	if !ok {
		return "(undefined)"
	}

	switch v := val.(type) {
	case string, int, float64, bool:
		return fmt.Sprintf("%v", v)
	case nil:
		return "(null)"
	default:
		j, err := json.Marshal(v)
		if err != nil {
			return "(invalid)"
		}
		return string(j)
	}
}

func (r *Repl) writeHistory(line string) error {
	configDirs := configdir.New(VENDOR_NAME, APP_NAME)
	folders := configDirs.QueryFolders(configdir.Global)
	if len(folders) == 0 {
		return fmt.Errorf("no config folder")
	}
	folder := folders[0]
	err := folder.MkdirAll()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(folder.Path, HISTORY_FILE), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(line + "\n"); err != nil {
		return err
	}
	return nil

}

func (r *Repl) readHistory() []string {
	configDirs := configdir.New(VENDOR_NAME, APP_NAME)
	folder := configDirs.QueryFolderContainsFile(HISTORY_FILE)
	if folder != nil {
		data, err := folder.ReadFile(HISTORY_FILE)
		if err != nil {
			return []string{}
		}
		return strings.Split(string(data), "\n")
	}
	return []string{}
}

func (r *Repl) pagerableOut() (io.Writer, func() error) {
	var buffer bytes.Buffer

	var out io.Writer = &buffer
	var pager *exec.Cmd
	if r.enabledPager {
		pager = exec.Command(getPagerCmd())
		pager.Stdin = &buffer
		pager.Stdout = r.out
	}

	if pager != nil {
		return out, func() error {
			if err := pager.Run(); err != nil {
				return err
			}
			return nil
		}
	}
	return out, func() error {
		fmt.Fprint(r.out, buffer.String())
		return nil
	}
}

func getPagerCmd() string {
	if env := os.Getenv("PAGER"); env != "" {
		return env
	}
	return "less"
}
