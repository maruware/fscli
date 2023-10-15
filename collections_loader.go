package fscli

import (
	"sync"

	"cloud.google.com/go/firestore"
)

var (
	baseDocToFetched     *sync.Map
	baseDocToCollections *sync.Map
)

func init() {
	baseDocToFetched = new(sync.Map)
	baseDocToCollections = new(sync.Map)
}

func shouldFetchCollections(baseDoc string) bool {
	_, ok := baseDocToFetched.Load(baseDoc)
	return !ok
}

func markCollectionsFetched(baseDoc string) {
	baseDocToFetched.Store(baseDoc, true)
}

func unmarkCollectionsFetched(baseDoc string) {
	baseDocToFetched.Delete(baseDoc)
}

func fetchCollections(baseDoc string, getCollectionsIterator func(baseDoc string) (*firestore.CollectionIterator, error)) {
	if !shouldFetchCollections(baseDoc) {
		return
	}

	markCollectionsFetched(baseDoc)
	itr, err := getCollectionsIterator(baseDoc)
	if err != nil {
		unmarkCollectionsFetched(baseDoc)
		return
	}

	collections := make([]string, 0)
	for {
		col, err := itr.Next()
		if err != nil {
			break
		}
		collections = append(collections, col.ID)

		if len(collections)%50 == 0 {
			baseDocToCollections.Store(baseDoc, collections)
		}
	}
	baseDocToCollections.Store(baseDoc, collections)
}

func getCollections(baseDoc string, getCollectionsIterator func(baseDoc string) (*firestore.CollectionIterator, error)) []string {
	go fetchCollections(baseDoc, getCollectionsIterator)
	if collection, ok := baseDocToCollections.Load(baseDoc); !ok {
		return []string{}
	} else {
		return collection.([]string)
	}
}
