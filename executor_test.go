package fscli

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func seed(c *firestore.Client, usersCollection string) error {
	ctx := context.Background()

	g := errgroup.Group{}
	for i := 0; i < 5; i++ {
		i := i
		g.Go(func() error {
			userId := strconv.Itoa(i)
			_, err := c.Collection(usersCollection).Doc(userId).Set(ctx, map[string]interface{}{
				"name": fmt.Sprintf("user-%d", i),
				"age":  20 + i,
				"nicknames": []string{
					fmt.Sprintf("u-%d-1", i),
					fmt.Sprintf("u-%d-2", i),
				},
			})
			if err != nil {
				return err
			}
			postId := fmt.Sprintf("post%d", i)
			postsCollection := fmt.Sprintf("%s/%s/posts", usersCollection, userId)
			_, err = c.Collection(postsCollection).Doc(postId).Set(ctx, map[string]interface{}{
				"title": fmt.Sprintf("post-%d", i),
			})
			if err != nil {
				return err
			}

			return nil
		})
	}

	err := g.Wait()
	return err
}

func cleanSeed(c *firestore.Client, usersCollection string) error {
	ctx := context.Background()

	usersItr := c.Collection(usersCollection).Query.Documents(ctx)
	for {
		doc, err := usersItr.Next()
		if err != nil {
			break
		}
		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			return err
		}

		postCollection := fmt.Sprintf("%s/%s/posts", usersCollection, doc.Ref.ID)
		postItr := doc.Ref.Collection(postCollection).Query.Documents(ctx)
		for {
			postDoc, err := postItr.Next()
			if err != nil {
				break
			}
			_, err = postDoc.Ref.Delete(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func TestQuery(t *testing.T) {
	os.Setenv("FIRESTORE_EMULATOR_HOST", "127.0.0.1:8080")
	ctx := context.Background()
	fs, err := firestore.NewClient(ctx, "test-project")
	if err != nil {
		t.Fatal(err)
	}
	exe := NewExecutor(ctx, fs)
	if err != nil {
		t.Fatal(err)
	}

	if exe == nil {
		t.Fatal("executor is nil")
	}

	prefix := "fscli-executor-test-query"
	usersCollection := prefix + "-users"

	err = seed(fs, usersCollection)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanSeed(fs, usersCollection)

	tests := []struct {
		desc  string
		input *QueryOperation
		want  []map[string]any
	}{
		{
			desc:  "simple query",
			input: NewQueryOperation(usersCollection, nil, nil),
			want: []map[string]any{
				{"name": "user-0", "age": int64(20), "nicknames": []any{"u-0-1", "u-0-2"}},
				{"name": "user-1", "age": int64(21), "nicknames": []any{"u-1-1", "u-1-2"}},
				{"name": "user-2", "age": int64(22), "nicknames": []any{"u-2-1", "u-2-2"}},
				{"name": "user-3", "age": int64(23), "nicknames": []any{"u-3-1", "u-3-2"}},
				{"name": "user-4", "age": int64(24), "nicknames": []any{"u-4-1", "u-4-2"}},
			},
		},
		{
			desc: "query with where",
			input: NewQueryOperation(usersCollection, nil, []Filter{
				NewStringFilter("name", "==", "user-1"),
			}),
			want: []map[string]any{
				{
					"name":      "user-1",
					"age":       int64(21),
					"nicknames": []any{"u-1-1", "u-1-2"},
				},
			},
		},
		{
			desc: "query with subcollection",
			input: NewQueryOperation(fmt.Sprintf("%s/1/posts", usersCollection), nil, []Filter{
				NewStringFilter("title", "==", "post-1"),
			}),
			want: []map[string]any{
				{
					"title": "post-1",
				},
			},
		},
		{
			desc: "query with not equal",
			input: NewQueryOperation(usersCollection, nil, []Filter{
				NewStringFilter("name", "!=", "user-1"),
			}),
			want: []map[string]any{
				{"name": "user-0", "age": int64(20), "nicknames": []any{"u-0-1", "u-0-2"}},
				{"name": "user-2", "age": int64(22), "nicknames": []any{"u-2-1", "u-2-2"}},
				{"name": "user-3", "age": int64(23), "nicknames": []any{"u-3-1", "u-3-2"}},
				{"name": "user-4", "age": int64(24), "nicknames": []any{"u-4-1", "u-4-2"}},
			},
		},
		{
			desc: "query with and",
			input: NewQueryOperation(usersCollection, nil, []Filter{
				NewStringFilter("name", "==", "user-1"),
				NewIntFilter("age", "==", 21),
			}),
			want: []map[string]any{
				{"name": "user-1", "age": int64(21), "nicknames": []any{"u-1-1", "u-1-2"}},
			},
		},
		{
			desc: "query with IN",
			input: NewQueryOperation(usersCollection, nil, []Filter{
				NewArrayFilter("age", "in", []any{20, 21, 22}),
			}),
			want: []map[string]any{
				{"name": "user-0", "age": int64(20), "nicknames": []any{"u-0-1", "u-0-2"}},
				{"name": "user-1", "age": int64(21), "nicknames": []any{"u-1-1", "u-1-2"}},
				{"name": "user-2", "age": int64(22), "nicknames": []any{"u-2-1", "u-2-2"}},
			},
		},
		{
			desc: "query with array-contains",
			input: NewQueryOperation(usersCollection, nil, []Filter{
				NewStringFilter("nicknames", "array-contains", "u-1-1"),
			}),
			want: []map[string]any{
				{"name": "user-1", "age": int64(21), "nicknames": []any{"u-1-1", "u-1-2"}},
			},
		},
		{
			desc: "query with array-contains-any",
			input: NewQueryOperation(usersCollection, nil, []Filter{
				NewArrayFilter("nicknames", "array-contains-any", []any{"u-1-1", "u-2-1"}),
			}),
			want: []map[string]any{
				{"name": "user-1", "age": int64(21), "nicknames": []any{"u-1-1", "u-1-2"}},
				{"name": "user-2", "age": int64(22), "nicknames": []any{"u-2-1", "u-2-2"}},
			},
		},
		{
			desc:  "query with select",
			input: NewQueryOperation(usersCollection, []string{"name", "age"}, nil),
			want: []map[string]any{
				{"name": "user-0", "age": int64(20)},
				{"name": "user-1", "age": int64(21)},
				{"name": "user-2", "age": int64(22)},
				{"name": "user-3", "age": int64(23)},
				{"name": "user-4", "age": int64(24)},
			},
		},
		{
			desc:  "query with select and where",
			input: NewQueryOperation(usersCollection, []string{"name", "age"}, []Filter{NewStringFilter("name", "==", "user-1")}),
			want: []map[string]any{
				{"name": "user-1", "age": int64(21)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			docs, err := exe.ExecuteQuery(ctx, tt.input)
			if err != nil {
				t.Fatal(err)
			}

			results := make([]map[string]any, 0)
			for _, doc := range docs {
				results = append(results, doc.Data())
			}
			assert.Equal(t, tt.want, results)
		})
	}
}

func TestGet(t *testing.T) {
	os.Setenv("FIRESTORE_EMULATOR_HOST", "127.0.0.1:8080")
	ctx := context.Background()
	fs, err := firestore.NewClient(ctx, "test-project")
	if err != nil {
		t.Fatal(err)
	}
	exe := NewExecutor(ctx, fs)
	if err != nil {
		t.Fatal(err)
	}

	if exe == nil {
		t.Fatal("executor is nil")
	}

	prefix := "fscli-executor-test-get"
	usersCollection := prefix + "-users"

	err = seed(fs, usersCollection)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanSeed(fs, usersCollection)

	tests := []struct {
		desc  string
		input *GetOperation
		want  map[string]any
	}{
		{
			desc:  "get",
			input: NewGetOperation(usersCollection, "1"),
			want: map[string]any{
				"name": "user-1", "age": int64(21), "nicknames": []any{"u-1-1", "u-1-2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			doc, err := exe.ExecuteGet(ctx, tt.input)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, doc.Data())
		})
	}
}
