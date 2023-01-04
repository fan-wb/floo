package core

import (
	"floo/catfs/db"
	n "floo/catfs/nodes"
	h "floo/util/hashlib"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"path"
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

// WithDummyLinker creates a testing linker and passes it to `fn`.
func WithDummyLinker(t *testing.T, fn func(lkr *Linker)) {
	WithDummyKv(t, func(kv db.Database) {
		lkr := NewLinker(kv)
		require.Nil(t, lkr.SetOwner("alice"))
		MustCommit(t, lkr, "init")

		fn(lkr)
	})
}

// MustCommit commits the current state with `msg`.
func MustCommit(t *testing.T, lkr *Linker, msg string) *n.Commit {
	if err := lkr.MakeCommit(n.AuthorOfStage, msg); err != nil {
		t.Fatalf("Failed to make commit with msg %s: %v", msg, err)
	}

	head, err := lkr.Head()
	if err != nil {
		t.Fatalf("Failed to retrieve head after commit: %v", err)
	}

	return head
}

// MustTouch creates a new node at `touchPath` and sets its content hash
// to a hash derived from `seed`.
func MustTouch(t *testing.T, lkr *Linker, touchPath string, seed byte) *n.File {
	dirname := path.Dir(touchPath)
	parent, err := lkr.LookupDirectory(dirname)
	if err != nil {
		t.Fatalf("touch: Failed to lookup: %s", dirname)
	}

	basePath := path.Base(touchPath)
	file := n.NewEmptyFile(parent, basePath, lkr.owner, lkr.NextInode())

	file.SetBackend(lkr, h.TestDummy(t, seed))
	file.SetContent(lkr, h.TestDummy(t, seed))
	file.SetKey(make([]byte, 32))

	child, err := parent.Child(lkr, basePath)
	if err != nil {
		t.Fatalf("touch: Failed to lookup child: %v %v", touchPath, err)
	}

	if child != nil {
		if err := parent.RemoveChild(lkr, child); err != nil {
			t.Fatalf("touch: failed to remove previous node: %v", err)
		}
	}

	if err := parent.Add(lkr, file); err != nil {
		t.Fatalf("touch: Adding %s to root failed: %v", touchPath, err)
	}

	if err := lkr.StageNode(file); err != nil {
		t.Fatalf("touch: Staging %s failed: %v", touchPath, err)
	}

	return file
}

// MustTouchAndCommit is a combined MustTouch and MustCommit.
func MustTouchAndCommit(t *testing.T, lkr *Linker, path string, seed byte) (*n.File, *n.Commit) {
	file, err := Stage(lkr, path, h.TestDummy(t, seed), h.TestDummy(t, seed), uint64(seed), nil)
	if err != nil {
		t.Fatalf("Failed to stage %s at %d: %v", path, seed, err)
	}

	return file, MustCommit(t, lkr, fmt.Sprintf("cmt %d", seed))
}
