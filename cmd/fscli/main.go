package main

import (
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/maruware/fscli"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project-id",
				Usage:    "firebase project id",
				Required: true,
			},
		},
		Action: func(cCtx *cli.Context) error {
			projectId := cCtx.String("project-id")
			fmt.Printf("project ID: %s\n", projectId)

			fs, err := firestore.NewClient(cCtx.Context, projectId)
			if err != nil {
				return err
			}
			defer fs.Close()

			fscli.ReplStart(cCtx.Context, fs, os.Stdin, os.Stdout)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
