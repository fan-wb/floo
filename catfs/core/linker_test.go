package core

import (
	"floo/catfs/db"
	n "floo/catfs/nodes"
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
