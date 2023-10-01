package fscli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"cloud.google.com/go/firestore"
)

const PROMPT = ">> "

func ReplStart(ctx context.Context, fs *firestore.Client, in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	executor := NewExecutor(ctx, fs)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if strings.TrimSpace(line) == "\\d" {
			listUpCollections(ctx, fs, out)
			continue
		}

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
