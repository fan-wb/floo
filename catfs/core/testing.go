package core

import (
	"floo/catfs/db"
	"os"
	"testing"
)

// WithDummyKv creates a testing key value store and passes it to `fn`.
func WithDummyKv(t *testing.T, fn func(kv db.Database)) {
	dbPath, err := os.MkdirTemp("", "floo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer os.RemoveAll(dbPath)

	kv, err := db.NewDiskDatabase(dbPath)
	if err != nil {
		t.Fatalf("Could not create dummy kv for tests: %v", err)
	}

	fn(kv)

	if err := kv.Close(); err != nil {
		t.Fatalf("Closing the dummy kv failed: %v", err)
	}
}
