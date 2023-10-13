package fscli

import (
	"context"
	"errors"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Executor struct {
	fs *firestore.Client
}

var ErrInvalidCollection = errors.New("invalid collection")

func NewExecutor(ctx context.Context, fs *firestore.Client) *Executor {
	return &Executor{fs}
}

func (exe *Executor) ExecuteQuery(ctx context.Context, op *QueryOperation) ([]*firestore.DocumentSnapshot, error) {
	collection := exe.fs.Collection(op.Collection())
	if collection == nil {
		return nil, ErrInvalidCollection
	}
	q := collection.Query
	for _, filter := range op.filters {
		q = q.Where(filter.FieldName(), string(filter.Operator()), filter.Value())
	}

	if len(op.selects) > 0 {
		q = q.Select(op.selects...)
	}

	if len(op.orderBys) > 0 {
		for _, orderBy := range op.orderBys {
			q = q.OrderBy(orderBy.field, firestore.Direction(orderBy.direction))
		}
	}

	if op.limit > 0 {
		q = q.Limit(op.limit)
	}

	itr := q.Documents(ctx)
	defer itr.Stop()

	docs := make([]*firestore.DocumentSnapshot, 0)
	for {
		doc, err := itr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		docs = append(docs, doc)
	}
	return docs, nil
}

func (exe *Executor) ExecuteGet(ctx context.Context, op *GetOperation) (*firestore.DocumentSnapshot, error) {
	doc, err := exe.fs.Collection(op.Collection()).Doc(op.DocId()).Get(ctx)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (exe *Executor) ExecuteListCollections(ctx context.Context, cmd *MetacommandListCollections) ([]string, error) {
	parts := strings.Split(cmd.baseDoc, "/")

	docParts := make([][2]string, 0)
	for i, part := range parts {
		if i%2 == 1 {
			docParts = append(docParts, [2]string{parts[i-1], part})
		}
	}

	var docBase *firestore.DocumentRef
	for _, docPart := range docParts {
		if docBase == nil {
			docBase = exe.fs.Collection(docPart[0]).Doc(docPart[1])
		} else {
			docBase = docBase.Collection(docPart[0]).Doc(docPart[1])
		}
	}

	if docBase == nil {
		return allCollections(exe.fs.Collections(ctx)), nil
	}

	return allCollections(docBase.Collections(ctx)), nil
}

func allCollections(itr *firestore.CollectionIterator) []string {
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
