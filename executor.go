package fscli

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Executor struct {
	c *firestore.Client
}

func NewExecutor(ctx context.Context, projectID string) (*Executor, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	e := &Executor{client}
	return e, nil
}

func (exe *Executor) ExecuteQuery(ctx context.Context, op *QueryOperation) ([]map[string]any, error) {
	collection := exe.c.Collection(op.Collection())
	q := collection.Query
	for _, filter := range op.filters {
		q = q.Where(filter.FieldName(), filter.Operator(), filter.Value())
	}

	itr := q.Documents(ctx)
	defer itr.Stop()

	results := []map[string]any{}
	for {
		doc, err := itr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		results = append(results, doc.Data())
	}
	return results, nil
}
