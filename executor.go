package fscli

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Executor struct {
	fs *firestore.Client
}

func NewExecutor(ctx context.Context, fs *firestore.Client) *Executor {
	return &Executor{fs}
}

func (exe *Executor) ExecuteQuery(ctx context.Context, op *QueryOperation) ([]*firestore.DocumentSnapshot, error) {
	collection := exe.fs.Collection(op.Collection())
	q := collection.Query
	for _, filter := range op.filters {
		q = q.Where(filter.FieldName(), string(filter.Operator()), filter.Value())
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
