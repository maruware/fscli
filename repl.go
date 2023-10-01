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
		op := parser.Parse()
		if op, ok := op.(*QueryOperation); ok {
			results, err := executor.ExecuteQuery(ctx, op)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				continue
			}

			for _, result := range results {
				j, err := json.Marshal(result)
				if err != nil {
					fmt.Printf("error: %s\n", err)
					continue
				}
				out.Write([]byte(fmt.Sprintf("%s\n", j)))
			}
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
