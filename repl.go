package fscli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/c-bata/go-prompt"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/exp/slices"
)

const LongLine = "--------------------------------------------------------------------------"

type OutputMode string

const (
	OutputModeJSON  OutputMode = "json"
	OutputModeTable OutputMode = "table"
)

type Repl struct {
	ctx        context.Context
	fs         *firestore.Client
	in         io.Reader
	out        io.Writer
	outputMode OutputMode
	exe        *Executor
}

func NewRepl(ctx context.Context, fs *firestore.Client, in io.Reader, out io.Writer, outputMode OutputMode) *Repl {
	return &Repl{
		ctx:        ctx,
		fs:         fs,
		in:         in,
		out:        out,
		outputMode: outputMode,
		exe:        NewExecutor(ctx, fs),
	}
}

func (r *Repl) completer(d prompt.Document) []prompt.Suggest {
	w := d.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}

	text := d.TextBeforeCursor()
	// trim inputting last word
	text = strings.TrimSuffix(text, d.GetWordBeforeCursor())

	c := NewCompleter(NewLexer(text))
	suggestions, err := c.Parse()
	if err != nil {
		return []prompt.Suggest{}
	}

	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func (r *Repl) Start() {
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
	)
	p.Run()
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

func (r *Repl) outputDocsTable(docs []*firestore.DocumentSnapshot) {
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

	table := tablewriter.NewWriter(r.out)
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

func (r *Repl) processLine(line string) {
	if strings.TrimSpace(line) == "" {
		return
	}
	if strings.TrimSpace(line) == "\\d" {
		listUpCollections(r.ctx, r.fs, r.out)
		return
	}

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
	if op, ok := op.(*QueryOperation); ok {
		docs, err := r.exe.ExecuteQuery(r.ctx, op)
		if err != nil {
			fmt.Fprintf(r.out, "error: %s\n", err)
			return
		}

		if r.outputMode == OutputModeJSON {
			for _, doc := range docs {
				r.outputDocJSON(doc)
				fmt.Fprintln(r.out, LongLine)
			}
		} else if r.outputMode == OutputModeTable {
			r.outputDocsTable(docs)
		}
	}
	if op, ok := op.(*GetOperation); ok {
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

func listUpCollections(ctx context.Context, fs *firestore.Client, out io.Writer) {
	cols := fs.Collections(ctx)
	for {
		col, err := cols.Next()
		if err != nil {
			break
		}
		fmt.Fprintf(out, "Collection: %s\n", col.ID)
	}
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
