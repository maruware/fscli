package fscli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/mattn/go-colorable"

	"github.com/fatih/color"
	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/coloring"
	"github.com/nyaosorg/go-readline-ny/simplehistory"
)

const LongLine = "--------------------------------------------------------------------------"

func ReplStart(ctx context.Context, fs *firestore.Client, in io.Reader, out io.Writer) {
	history := simplehistory.New()

	editor := &readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) {
			green := color.New(color.FgGreen)
			return green.Fprintf(w, "> ")
		},
		Writer:         colorable.NewColorableStdout(),
		History:        history,
		Coloring:       &coloring.VimBatch{},
		HistoryCycling: true,
	}

	executor := NewExecutor(ctx, fs)

	for {
		line, err := editor.ReadLine(ctx)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("error: %s\n", err)
			return
		}

		if strings.TrimSpace(line) == "\\d" {
			listUpCollections(ctx, fs, out)
			continue
		}

		history.Add(line)

		lexer := NewLexer(line)
		parser := NewParser(lexer)
		op, err := parser.Parse()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
		if op == nil {
			continue
		}
		if op, ok := op.(*QueryOperation); ok {
			docs, err := executor.ExecuteQuery(ctx, op)
			if err != nil {
				fmt.Fprintf(out, "error: %s\n", err)
				continue
			}

			for _, doc := range docs {
				fmt.Fprintf(out, "ID: %s\n", doc.Ref.ID)

				j, err := json.Marshal(doc.Data())
				if err != nil {
					fmt.Printf("invalid data: %s\n", err)
					fmt.Fprintln(out, LongLine)
					continue
				}
				fmt.Fprintf(out, "Data: %s\n", j)
				fmt.Fprintln(out, LongLine)
			}
		}
		if op, ok := op.(*GetOperation); ok {
			result, err := executor.ExecuteGet(ctx, op)
			if err != nil {
				fmt.Fprintf(out, "error: %s\n", err)
				continue
			}
			j, err := json.Marshal(result)
			if err != nil {
				fmt.Fprintf(out, "error: %s\n", err)
				continue
			}
			fmt.Fprintf(out, "%s\n", j)
		}
	}
}

func listUpCollections(ctx context.Context, fs *firestore.Client, out io.Writer) {
	cols := fs.Collections(ctx)
	for {
		col, err := cols.Next()
		if err != nil {
			break
		}
		out.Write([]byte(fmt.Sprintf("%s\n", col.ID)))
	}
}
