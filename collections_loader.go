package fscli

import (
	"log"
	"sync"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

var (
	baseDocToFetched     *sync.Map
	baseDocToCollections *sync.Map
)

const PAGE_SIZE = 100

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
	p := iterator.NewPager(itr, PAGE_SIZE, "")
	for {
		var cols []*firestore.CollectionRef
		pageToken, err := p.NextPage(&cols)
		if err != nil {
			log.Fatal(err)
			return
		}
		collections = append(collections, getCollectionIds(cols)...)
		baseDocToCollections.Store(baseDoc, collections)

		if pageToken == "" {
			break
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
