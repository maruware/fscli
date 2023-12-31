package fscli

import (
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
		r.processLine,
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

func (r *Repl) processLine(line string) {
	if strings.TrimSpace(line) == "" {
		return
	}

	err := r.writeHistory(line)
	if err != nil {
		fmt.Fprintf(r.out, "error: %s\n", err)
		return
	}

	lexer := NewLexer(line)
	parser := NewParser(lexer)
	result, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(r.out, "error: %s\n", err)
		return
	}
	if result == nil {
		return
	}

	if op, ok := result.(*MetacommandPager); ok {
		r.enabledPager = op.on
		return
	}

	if op, ok := result.(*MetacommandListCollections); ok {
		cols, err := r.exe.ExecuteListCollections(r.ctx, op)
		if err != nil {
			fmt.Fprintf(r.out, "error: %s\n", err)
			return
		}

		out, render := r.pagerableOut()
		for _, col := range cols {
			fmt.Fprintf(out, "%s\n", col)
		}
		if err := render(); err != nil {
			fmt.Fprintf(r.out, "error: %s\n", err)
		}
	}

	if op, ok := result.(*QueryOperation); ok {
		docs, err := r.exe.ExecuteQuery(r.ctx, op)
		if err != nil {
			fmt.Fprintf(r.out, "error: %s\n", err)
			return
		}

		if r.outputMode == OutputModeJSON {
			r.oututDocsJSON(docs)
		} else if r.outputMode == OutputModeTable {
			r.outputDocsTable(docs)
		}
	}
	if op, ok := result.(*GetOperation); ok {
		doc, err := r.exe.ExecuteGet(r.ctx, op)
		if err != nil {
			fmt.Fprintf(r.out, "error: %s\n", err)
			return
		}

		if r.outputMode == OutputModeJSON {
			r.outputDocJSON(doc)
		} else if r.outputMode == OutputModeTable {
			r.outputDocTable(doc)
		}
	}
	if op, ok := result.(*CountOperation); ok {
		count, err := r.exe.ExecuteCount(r.ctx, op)
		if err != nil {
			fmt.Fprintf(r.out, "error: %s\n", err)
			return
		}

		fmt.Fprintf(r.out, "%d\n", count)
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

func (r *Repl) oututDocsJSON(docs []*firestore.DocumentSnapshot) {
	for _, doc := range docs {
		r.outputDocJSON(doc)
		fmt.Fprintln(r.out, LongLine)
	}
}

func (r *Repl) outputDocJSON(doc *firestore.DocumentSnapshot) {
	fmt.Fprintf(r.out, "ID: %s\n", doc.Ref.ID)
	j, err := json.Marshal(doc.Data())
	if err != nil {
		fmt.Fprintf(r.out, "invalid data: %s\n", err)
		return
	}
	fmt.Fprintf(r.out, "Data: %s\n", j)
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
