package core

import (
	"floo/catfs/db"
	ie "floo/catfs/errors"
	n "floo/catfs/nodes"
	h "floo/util/hashlib"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"unsafe"
)

// Basic test to see if the root node can be inserted and stored.
// A new staging commit should be also created in the background.
// On the second run, the root node should be already cached.
func TestLinkerInsertRoot(t *testing.T) {
	WithDummyKv(t, func(kv db.Database) {
		lkr := NewLinker(kv)
		root, err := n.NewEmptyDirectory(lkr, nil, "/", "u", 2)
		if err != nil {
			t.Fatalf("Creating empty root dir failed: %v", err)
		}

		if err := lkr.StageNode(root); err != nil {
			t.Fatalf("Staging root failed: %v", err)
		}

		sameRoot, err := lkr.ResolveDirectory("/")
		if err != nil {
			t.Fatalf("Resolving root failed: %v", err)
		}

		if sameRoot == nil {
			t.Fatal("Resolving root failed (is nil)")
		}

		if path := sameRoot.Path(); path != "/" {
			t.Fatalf("Path of root is not /: %s", path)
		}

		ptrRoot, err := lkr.ResolveDirectory("/")
		if err != nil {
			t.Fatalf("Second lookup of root failed: %v", err)
		}

		if unsafe.Pointer(ptrRoot) != unsafe.Pointer(sameRoot) {
			t.Fatal("Second root did not come from the cache")
		}

		status, err := lkr.Status()
		if err != nil {
			t.Fatalf("Failed to retrieve status: %v", err)
		}

		if !status.Root().Equal(root.TreeHash()) {
			t.Fatalf("status.root and root differ: %v <-> %v", status.Root(), root.TreeHash())
		}
	})
}

func TestLinkerRefs(t *testing.T) {
	author := n.AuthorOfStage
	WithDummyKv(t, func(kv db.Database) {
		lkr := NewLinker(kv)
		root, err := lkr.Root()
		if err != nil {
			t.Fatalf("Failed to create root: %v", err)
		}

		newFile := n.NewEmptyFile(root, "cat.png", "u", 2)
		if err != nil {
			t.Fatalf("Failed to create empty file: %v", err)
		}

		newFile.SetSize(10)
		newFile.SetContent(lkr, h.TestDummy(t, 1))

		if err := root.Add(lkr, newFile); err != nil {
			t.Fatalf("Adding empty file failed: %v", err)
		}

		if err := lkr.StageNode(newFile); err != nil {
			t.Fatalf("Staging new file failed: %v", err)
		}

		if _, err := lkr.Head(); !ie.IsErrNoSuchRef(err) {
			t.Fatalf("There is a HEAD from start?!")
		}

		cmt, err := lkr.Status()
		if err != nil || cmt == nil {
			t.Fatalf("Failed to retrieve status: %v", err)
		}

		if err := lkr.MakeCommit(author, "First commit"); err != nil {
			t.Fatalf("Making commit failed: %v", err)
		}

		// Assert that staging is empty (except the "/stage/STATUS" part)
		foundKeys := []string{}
		keys, err := kv.Keys("stage")
		require.Nil(t, err)

		for _, key := range keys {
			foundKeys = append(foundKeys, strings.Join(key, "/"))
		}

		require.Equal(t, []string{"stage/STATUS"}, foundKeys)

		head, err := lkr.Head()
		if err != nil {
			t.Fatalf("Obtaining HEAD failed: %v", err)
		}

		status, err := lkr.Status()
		if err != nil {
			t.Fatalf("Failed to obtain the status: %v", err)
		}

		if !head.Root().Equal(status.Root()) {
			t.Fatalf("HEAD and CURR are not equal after first commit.")
		}

		if err := lkr.MakeCommit(author, "No."); err != ie.ErrNoChange {
			t.Fatalf("Committing without change led to a new commit: %v", err)
		}
	})
}

// Test if Linker can load objects after closing/re-opening the kv.
func TestLinkerPersistence(t *testing.T) {
	dbPath, err := os.MkdirTemp("", "floo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer os.RemoveAll(dbPath)

	kv, err := db.NewDiskDatabase(dbPath)
	if err != nil {
		t.Fatalf("Could not create dummy kv for tests: %v", err)
	}

	lkr := NewLinker(kv)
	if err := lkr.MakeCommit(n.AuthorOfStage, "initial commit"); err != nil {
		t.Fatalf("Failed to create initial commit out of nothing: %v", err)
	}

	head, err := lkr.Head()
	if err != nil {
		t.Fatalf("Failed to retrieve Head after initial commit: %v", err)
	}

	oldHeadHash := head.TreeHash().Clone()

	if err := kv.Close(); err != nil {
		t.Fatalf("Closing the dummy kv failed: %v", err)
	}

	kv, err = db.NewDiskDatabase(dbPath)
	if err != nil {
		t.Fatalf("Could not create second dummy kv: %v", err)
	}

	lkr = NewLinker(kv)
	head, err = lkr.Head()
	if err != nil {
		t.Fatalf("Failed to retrieve head after kv reload: %v", err)
	}

	if !oldHeadHash.Equal(head.TreeHash()) {
		t.Fatalf("HEAD hash differs before and after reload: %v <-> %v", oldHeadHash, head.TreeHash())
	}

	if err := kv.Close(); err != nil {
		t.Fatalf("Closing the second kv failed: %v", err)
	}
}
