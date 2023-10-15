package fscli

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
)

func findAllCollections(ctx context.Context, fs *firestore.Client, baseDoc string) ([]string, error) {
	itr, err := getCollectionsIterator(ctx, fs, baseDoc)
	if err != nil {
		return nil, err
	}
	return iterateAllCollections(itr), nil
}

func getCollectionsIterator(ctx context.Context, fs *firestore.Client, baseDoc string) (*firestore.CollectionIterator, error) {
	if baseDoc == "" {
		return fs.Collections(ctx), nil
	}

	lastSlash := strings.LastIndex(baseDoc, "/")
	if lastSlash == -1 {
		return nil, fmt.Errorf("invalid path: %s", baseDoc)
	}

	collection := baseDoc[:lastSlash]
	docId := baseDoc[lastSlash+1:]

	return fs.Collection(collection).Doc(docId).Collections(ctx), nil
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
