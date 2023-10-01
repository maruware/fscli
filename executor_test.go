package fscli

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
)

func seed(c *firestore.Client) error {
	ctx := context.Background()
	_, err := c.Collection("users").Doc("user1").Set(ctx, map[string]interface{}{
		"name": "user1",
		"age":  20,
	})
	if err != nil {
		return err
	}
	_, err = c.Collection("users").Doc("user2").Set(ctx, map[string]interface{}{
		"name": "user2",
		"age":  30,
	})
	if err != nil {
		return err
	}
	_, err = c.Collection("users").Doc("user3").Set(ctx, map[string]interface{}{
		"name": "user3",
		"age":  40,
	})
	if err != nil {
		return err
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

	tests := []struct {
		desc  string
		input *QueryOperation
		want  []map[string]any
	}{
		{
			desc: "simple query",
			input: NewQueryOperation("users", []Filter{
				&StringFilter{
					field:    "name",
					operator: "==",
					value:    "user1",
				},
			}),
			want: []map[string]any{
				{
					"name": "user1",
					"age":  int64(20),
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
