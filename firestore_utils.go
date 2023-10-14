package fscli

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
)

func findAllCollections(ctx context.Context, fs *firestore.Client, baseDoc string) ([]string, error) {
	if baseDoc == "" {
		return iterateAllCollections(fs.Collections(ctx)), nil
	}

	lastSlash := strings.LastIndex(baseDoc, "/")
	if lastSlash == -1 {
		return nil, fmt.Errorf("invalid path: %s", baseDoc)
	}

	collection := baseDoc[:lastSlash]
	docId := baseDoc[lastSlash+1:]

	return iterateAllCollections(fs.Collection(collection).Doc(docId).Collections(ctx)), nil
}

func iterateAllCollections(itr *firestore.CollectionIterator) []string {
	cols := make([]string, 0)
	for {
		col, err := itr.Next()
		if err != nil {
			break
		}
		cols = append(cols, col.ID)
	}
	return cols
}
