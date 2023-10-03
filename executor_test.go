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

const prefix = "fscli-executor-test"
const usersCollection = prefix + "-users"

func seed(c *firestore.Client) error {
	ctx := context.Background()

	g := errgroup.Group{}
	for i := 0; i < 5; i++ {
		i := i
		g.Go(func() error {
			userId := strconv.Itoa(i)
			_, err := c.Collection(usersCollection).Doc(userId).Set(ctx, map[string]interface{}{
				"name": fmt.Sprintf("user-%d", i),
				"age":  20,
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

func cleanSeed(c *firestore.Client) error {
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
			desc: "simple query",
			input: NewQueryOperation(usersCollection, []Filter{
				NewStringFilter("name", "==", "user-1"),
			}),
			want: []map[string]any{
				{
					"name": "user-1",
					"age":  int64(20),
				},
			},
		},
		{
			desc: "simple query with subcollection",
			input: NewQueryOperation(fmt.Sprintf("%s/1/posts", usersCollection), []Filter{
				NewStringFilter("title", "==", "post-1"),
			}),
			want: []map[string]any{
				{
					"title": "post-1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			results, err := exe.ExecuteQuery(ctx, tt.input)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, results)
		})
	}
}
