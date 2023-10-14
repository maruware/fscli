package fscli

import "sync"

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

func fetchCollections(baseDoc string, findAllCollections func(baseDoc string) ([]string, error)) {
	if !shouldFetchCollections(baseDoc) {
		return
	}

	markCollectionsFetched(baseDoc)
	collection, err := findAllCollections(baseDoc)
	if err != nil {
		unmarkCollectionsFetched(baseDoc)
		return
	}
	baseDocToCollections.Store(baseDoc, collection)
}

func getCollections(baseDoc string, findAllCollections func(baseDoc string) ([]string, error)) []string {
	go fetchCollections(baseDoc, findAllCollections)
	if collection, ok := baseDocToCollections.Load(baseDoc); !ok {
		return []string{}
	} else {
		return collection.([]string)
	}
}
