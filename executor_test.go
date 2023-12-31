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

func seed(c *firestore.Client) error {
	ctx := context.Background()

	g := errgroup.Group{}
	for i := 0; i < 5; i++ {
		i := i
		g.Go(func() error {
			userId := strconv.Itoa(i)
			_, err := c.Collection("users").Doc(userId).Set(ctx, map[string]interface{}{
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
			postsCollection := fmt.Sprintf("%s/%s/posts", "users", userId)
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

func cleanSeed(c *firestore.Client) error {
	ctx := context.Background()

	usersItr := c.Collection("users").Query.Documents(ctx)
	for {
		doc, err := usersItr.Next()
		if err != nil {
			break
		}
		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			return err
		}

		postCollection := fmt.Sprintf("%s/%s/posts", "users", doc.Ref.ID)
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
	fs, err := firestore.NewClient(ctx, "fscli-executor-test-query")
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

	err = seed(fs)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanSeed(fs)

	tests := []struct {
		desc  string
		input *QueryOperation
		want  []map[string]any
	}{
		{
			desc:  "simple query",
			input: &QueryOperation{collection: "users"},
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
			input: &QueryOperation{collection: "users", filters: []Filter{
				NewStringFilter("name", "==", "user-1"),
			}},
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
			input: &QueryOperation{collection: fmt.Sprintf("%s/1/posts", "users"), filters: []Filter{
				NewStringFilter("title", "==", "post-1"),
			}},
			want: []map[string]any{
				{
					"title": "post-1",
				},
			},
		},
		{
			desc: "query with not equal",
			input: &QueryOperation{collection: "users", filters: []Filter{
				NewStringFilter("name", "!=", "user-1"),
			}},
			want: []map[string]any{
				{"name": "user-0", "age": int64(20), "nicknames": []any{"u-0-1", "u-0-2"}},
				{"name": "user-2", "age": int64(22), "nicknames": []any{"u-2-1", "u-2-2"}},
				{"name": "user-3", "age": int64(23), "nicknames": []any{"u-3-1", "u-3-2"}},
				{"name": "user-4", "age": int64(24), "nicknames": []any{"u-4-1", "u-4-2"}},
			},
		},
		{
			desc: "query with greater than",
			input: &QueryOperation{collection: "users", filters: []Filter{
				NewIntFilter("age", ">", 22),
			}},
			want: []map[string]any{
				{"name": "user-3", "age": int64(23), "nicknames": []any{"u-3-1", "u-3-2"}},
				{"name": "user-4", "age": int64(24), "nicknames": []any{"u-4-1", "u-4-2"}},
			},
		},
		{
			desc: "query with greater than or equal",
			input: &QueryOperation{collection: "users", filters: []Filter{
				NewIntFilter("age", ">=", 22),
			}},
			want: []map[string]any{
				{"name": "user-2", "age": int64(22), "nicknames": []any{"u-2-1", "u-2-2"}},
				{"name": "user-3", "age": int64(23), "nicknames": []any{"u-3-1", "u-3-2"}},
				{"name": "user-4", "age": int64(24), "nicknames": []any{"u-4-1", "u-4-2"}},
			},
		},
		{
			desc: "query with less than",
			input: &QueryOperation{collection: "users", filters: []Filter{
				NewIntFilter("age", "<", 22),
			}},
			want: []map[string]any{
				{"name": "user-0", "age": int64(20), "nicknames": []any{"u-0-1", "u-0-2"}},
				{"name": "user-1", "age": int64(21), "nicknames": []any{"u-1-1", "u-1-2"}},
			},
		},
		{
			desc: "query with less than or equal",
			input: &QueryOperation{collection: "users", filters: []Filter{
				NewIntFilter("age", "<=", 22),
			}},
			want: []map[string]any{
				{"name": "user-0", "age": int64(20), "nicknames": []any{"u-0-1", "u-0-2"}},
				{"name": "user-1", "age": int64(21), "nicknames": []any{"u-1-1", "u-1-2"}},
				{"name": "user-2", "age": int64(22), "nicknames": []any{"u-2-1", "u-2-2"}},
			},
		},
		{
			desc: "query with and",
			input: &QueryOperation{collection: "users", filters: []Filter{
				NewStringFilter("name", "==", "user-1"),
				NewIntFilter("age", "==", 21),
			}},
			want: []map[string]any{
				{"name": "user-1", "age": int64(21), "nicknames": []any{"u-1-1", "u-1-2"}},
			},
		},
		{
			desc: "query with IN",
			input: &QueryOperation{collection: "users", filters: []Filter{
				NewArrayFilter("age", "in", []any{20, 21, 22}),
			}},
			want: []map[string]any{
				{"name": "user-0", "age": int64(20), "nicknames": []any{"u-0-1", "u-0-2"}},
				{"name": "user-1", "age": int64(21), "nicknames": []any{"u-1-1", "u-1-2"}},
				{"name": "user-2", "age": int64(22), "nicknames": []any{"u-2-1", "u-2-2"}},
			},
		},
		{
			desc: "query with array-contains",
			input: &QueryOperation{collection: "users", filters: []Filter{
				NewStringFilter("nicknames", "array-contains", "u-1-1"),
			}},
			want: []map[string]any{
				{"name": "user-1", "age": int64(21), "nicknames": []any{"u-1-1", "u-1-2"}},
			},
		},
		{
			desc: "query with array-contains-any",
			input: &QueryOperation{collection: "users", filters: []Filter{
				NewArrayFilter("nicknames", "array-contains-any", []any{"u-1-1", "u-2-1"}),
			}},
			want: []map[string]any{
				{"name": "user-1", "age": int64(21), "nicknames": []any{"u-1-1", "u-1-2"}},
				{"name": "user-2", "age": int64(22), "nicknames": []any{"u-2-1", "u-2-2"}},
			},
		},
		{
			desc:  "query with select",
			input: &QueryOperation{collection: "users", selects: []string{"name", "age"}},
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
			input: &QueryOperation{collection: "users", selects: []string{"name", "age"}, filters: []Filter{NewStringFilter("name", "==", "user-1")}},
			want: []map[string]any{
				{"name": "user-1", "age": int64(21)},
			},
		},
		{
			desc:  "query with order by",
			input: &QueryOperation{collection: "users", orderBys: []OrderBy{{"age", firestore.Desc}}},
			want: []map[string]any{
				{"name": "user-4", "age": int64(24), "nicknames": []any{"u-4-1", "u-4-2"}},
				{"name": "user-3", "age": int64(23), "nicknames": []any{"u-3-1", "u-3-2"}},
				{"name": "user-2", "age": int64(22), "nicknames": []any{"u-2-1", "u-2-2"}},
				{"name": "user-1", "age": int64(21), "nicknames": []any{"u-1-1", "u-1-2"}},
				{"name": "user-0", "age": int64(20), "nicknames": []any{"u-0-1", "u-0-2"}},
			},
		},
		{
			desc:  "query with limit",
			input: &QueryOperation{collection: "users", limit: 2},
			want: []map[string]any{
				{"name": "user-0", "age": int64(20), "nicknames": []any{"u-0-1", "u-0-2"}},
				{"name": "user-1", "age": int64(21), "nicknames": []any{"u-1-1", "u-1-2"}},
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

func TestInvalidQuery(t *testing.T) {
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

	err = seed(fs)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanSeed(fs)

	tests := []struct {
		desc  string
		input *QueryOperation
		err   error
	}{
		{
			desc:  "query with invalid collection",
			input: &QueryOperation{collection: "users" + "/abc"},
			err:   ErrInvalidCollection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := exe.ExecuteQuery(ctx, tt.input)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestGet(t *testing.T) {
	os.Setenv("FIRESTORE_EMULATOR_HOST", "127.0.0.1:8080")
	ctx := context.Background()
	fs, err := firestore.NewClient(ctx, "fscli-executor-test-get")
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

	err = seed(fs)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanSeed(fs)

	tests := []struct {
		desc  string
		input *GetOperation
		want  map[string]any
	}{
		{
			desc:  "get",
			input: NewGetOperation("users", "1"),
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

func TestCount(t *testing.T) {
	os.Setenv("FIRESTORE_EMULATOR_HOST", "127.0.0.1:8080")
	ctx := context.Background()
	fs, err := firestore.NewClient(ctx, "fscli-executor-test-count")
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

	err = seed(fs)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanSeed(fs)

	tests := []struct {
		desc  string
		input *CountOperation
		want  int64
	}{
		{
			desc:  "count",
			input: NewCountOperation("users", []Filter{}),
			want:  5,
		},
		{
			desc: "count with where",
			input: NewCountOperation("users", []Filter{
				NewIntFilter("age", ">=", 22),
			}),
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			count, err := exe.ExecuteCount(ctx, tt.input)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, count)
		})
	}
}

func TestListCollections(t *testing.T) {
	seed := func(c *firestore.Client) error {
		ctx := context.Background()

		_, err := c.Collection("users").Doc("1").Set(ctx, map[string]interface{}{
			"name": "user-1",
		})
		if err != nil {
			return err
		}

		_, err = c.Collection("users").Doc("1").Collection("posts").Doc("1").Set(ctx, map[string]interface{}{
			"title": "post-1",
		})

		return err
	}

	usersCollection := "users"

	ctx := context.Background()
	fs, err := firestore.NewClient(ctx, "fscli-executor-test-list-cols")
	if err != nil {
		t.Fatal(err)
	}
	exe := NewExecutor(ctx, fs)
	if err != nil {
		t.Fatal(err)
	}

	seed(fs)

	tests := []struct {
		desc  string
		input *MetacommandListCollections
		want  []string
	}{
		{
			"list collections",
			&MetacommandListCollections{},
			[]string{usersCollection},
		},
		{
			desc:  "list subcollections",
			input: &MetacommandListCollections{baseDoc: fmt.Sprintf("%s/1", usersCollection)},
			want:  []string{"posts"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cols, err := exe.ExecuteListCollections(ctx, tt.input)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, cols)
		})
	}
}
