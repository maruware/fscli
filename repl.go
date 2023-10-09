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

var rootSuggestions = []prompt.Suggest{
	{Text: "GET", Description: "GET [docPath]"},
	{Text: "QUERY", Description: "QUERY [collection]"},
}

var querySuggestions = []prompt.Suggest{
	{Text: "SELECT", Description: "SELECT [field...]"},
	{Text: "WHERE", Description: "WHERE [field] [operator] [value]"},
}

func completer(d prompt.Document) []prompt.Suggest {
	w := d.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	if strings.HasPrefix(strings.ToUpper(d.CurrentLine()), "GET") {
		return []prompt.Suggest{}
	}
	// TODO: strict parse
	if strings.HasPrefix(strings.ToUpper(d.CurrentLine()), "QUERY") {
		return prompt.FilterHasPrefix(querySuggestions, d.GetWordBeforeCursor(), true)
	}

	return prompt.FilterHasPrefix(rootSuggestions, d.GetWordBeforeCursor(), true)
}

func (r *Repl) Start() {
	p := prompt.New(
		r.processLine,
		completer,
		prompt.OptionPrefix("> "),
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
	keys := []string{"ID"}
	for _, doc := range docs {
		for k := range doc.Data() {
			if !slices.Contains(keys, k) {
				keys = append(keys, k)
			}
		}
	}

	table := tablewriter.NewWriter(r.out)
	table.SetAutoFormatHeaders(false)
	table.SetHeader(keys)

	for _, doc := range docs {
		row := []string{}
		for _, k := range keys {
			if k == "ID" {
				row = append(row, doc.Ref.ID)
				continue
			}
			v, ok := doc.Data()[k]
			if !ok {
				row = append(row, "")
				continue
			}
			row = append(row, fmt.Sprintf("%v", v))
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
	keys := []string{"ID"}
	for k := range doc.Data() {
		keys = append(keys, k)
	}

	table := tablewriter.NewWriter(r.out)
	table.SetAutoFormatHeaders(false)
	table.SetHeader(keys)

	row := []string{}
	for _, k := range keys {
		if k == "ID" {
			row = append(row, doc.Ref.ID)
			continue
		}
		v, ok := doc.Data()[k]
		if !ok {
			row = append(row, "")
			continue
		}
		row = append(row, fmt.Sprintf("%v", v))
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
