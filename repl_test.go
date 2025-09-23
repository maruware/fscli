package fscli

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
)

func TestRepl_ProcessLineFromPipe(t *testing.T) {
	os.Setenv("FIRESTORE_EMULATOR_HOST", "127.0.0.1:8080")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fs, err := firestore.NewClient(ctx, "fscli-repl-test")
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()

	// Seed data
	docRef := fs.Collection("users").Doc("testuser")
	_, err = docRef.Set(ctx, map[string]interface{}{
		"name": "test user",
		"age":  30,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer docRef.Delete(ctx)

	var stdin bytes.Buffer
	var stdout bytes.Buffer

	stdin.Write([]byte("GET users/testuser\n"))

	repl := NewRepl(ctx, fs, &stdin, &stdout, OutputModeJSON)
	repl.ProcessLineFromPipe()

	expectedJSON := `{"id":"testuser","data":{"age":30,"name":"test user"}}`
	assert.JSONEq(t, expectedJSON, stdout.String())
}
