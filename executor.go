package fscli

import (
	"context"
	"errors"

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
	return findAllCollections(ctx, exe.fs, cmd.baseDoc)
}
