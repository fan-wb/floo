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

func TestCollideSameObjectHash(t *testing.T) {
	WithDummyKv(t, func(kv db.Database) {
		lkr := NewLinker(kv)
		root, err := lkr.Root()
		if err != nil {
			t.Fatalf("Failed to retrieve root: %v", err)
		}

		sub, err := n.NewEmptyDirectory(lkr, root, "sub", "u", 3)
		if err != nil {
			t.Fatalf("Creating empty sub dir failed: %v", err)
			return
		}

		if err := lkr.StageNode(sub); err != nil {
			t.Fatalf("Staging /sub failed: %v", err)
		}

		file1 := n.NewEmptyFile(sub, "a.png", "u", 4)
		if err != nil {
			t.Fatalf("Failed to create empty file1: %v", err)
		}

		file2 := n.NewEmptyFile(root, "a.png", "u", 5)
		if err != nil {
			t.Fatalf("Failed to create empty file2: %v", err)
		}

		file3 := n.NewEmptyFile(root, "b.png", "u", 6)
		if err != nil {
			t.Fatalf("Failed to create empty file3: %v", err)
		}

		file1.SetContent(lkr, h.TestDummy(t, 1))
		file2.SetContent(lkr, h.TestDummy(t, 1))
		file3.SetContent(lkr, h.TestDummy(t, 1))

		if err := sub.Add(lkr, file1); err != nil {
			t.Fatalf("Failed to add file1: %v", err)
		}
		if err := root.Add(lkr, file2); err != nil {
			t.Fatalf("Failed to add file2: %v", err)
		}
		if err := root.Add(lkr, file3); err != nil {
			t.Fatalf("Failed to add file3: %v", err)
		}

		if err := lkr.StageNode(file1); err != nil {
			t.Fatalf("Failed to stage file1: %v", err)
		}
		if err := lkr.StageNode(file2); err != nil {
			t.Fatalf("Failed to stage file2: %v", err)
		}
		if err := lkr.StageNode(file3); err != nil {
			t.Fatalf("Failed to stage file3: %v", err)
		}

		if file1.TreeHash().Equal(file2.TreeHash()) {
			t.Fatalf("file1 and file2 hash is equal: %v", file1.TreeHash())
		}
		if file2.TreeHash().Equal(file3.TreeHash()) {
			t.Fatalf("file2 and file3 hash is equal: %v", file2.TreeHash())
		}

		// Make sure we load the actual hashes from disk:
		lkr.MemIndexClear()
		file1Reset, err := lkr.LookupFile("/sub/a.png")
		if err != nil {
			t.Fatalf("Re-Lookup of file1 failed: %v", err)
		}
		file2Reset, err := lkr.LookupFile("/a.png")
		if err != nil {
			t.Fatalf("Re-Lookup of file2 failed: %v", err)
		}
		file3Reset, err := lkr.LookupFile("/b.png")
		if err != nil {
			t.Fatalf("Re-Lookup of file3 failed: %v", err)
		}

		if file1Reset.TreeHash().Equal(file2Reset.TreeHash()) {
			t.Fatalf("file1Reset and file2Reset hash is equal: %v", file1.TreeHash())
		}
		if file2Reset.TreeHash().Equal(file3Reset.TreeHash()) {
			t.Fatalf("file2Reset and file3Reset hash is equal: %v", file2.TreeHash())
		}
	})
}
func TestHaveStagedChanges(t *testing.T) {
	WithDummyLinker(t, func(lkr *Linker) {
		hasChanges, err := lkr.HaveStagedChanges()
		if err != nil {
			t.Fatalf("have staged changes failed before touch: %v", err)
		}
		if hasChanges {
			t.Fatalf("HaveStagedChanges has changes before something happened")
		}

		MustTouch(t, lkr, "/x.png", 1)

		hasChanges, err = lkr.HaveStagedChanges()
		if err != nil {
			t.Fatalf("have staged changes failed after touch: %v", err)
		}
		if !hasChanges {
			t.Fatalf("HaveStagedChanges has no changes after something happened")
		}

		MustCommit(t, lkr, "second")

		hasChanges, err = lkr.HaveStagedChanges()
		if err != nil {
			t.Fatalf("have staged changes failed after commit: %v", err)
		}
		if hasChanges {
			t.Fatalf("HaveStagedChanges has changes after commit")
		}
	})
}

func TestFilesByContent(t *testing.T) {
	WithDummyLinker(t, func(lkr *Linker) {
		file := MustTouch(t, lkr, "/x.png", 1)

		contents := []h.Hash{file.BackendHash()}
		result, err := lkr.FilesByContents(contents)

		require.Nil(t, err)

		resultFile, ok := result[file.BackendHash().B58String()]
		require.True(t, ok)
		require.Len(t, result, 1)
		require.Equal(t, file, resultFile)
	})
}

func TestResolveRef(t *testing.T) {
	WithDummyLinker(t, func(lkr *Linker) {
		initCmt, err := lkr.Head()
		require.Nil(t, err)

		cmts := []*n.Commit{initCmt}
		for idx := 0; idx < 10; idx++ {
			_, cmt := MustTouchAndCommit(t, lkr, "/x", byte(idx))
			cmts = append([]*n.Commit{cmt}, cmts...)
		}

		// Insert the init cmt a few times as fodder:
		cmts = append(cmts, initCmt)
		cmts = append(cmts, initCmt)
		cmts = append(cmts, initCmt)

		for nUp := 0; nUp < len(cmts)+3; nUp++ {
			refname := "head"
			for idx := 0; idx < nUp; idx++ {
				refname += "^"
			}

			expect := initCmt
			if nUp < len(cmts) {
				expect = cmts[nUp]
			}

			ref, err := lkr.ResolveRef(refname)
			require.Nil(t, err)
			require.Equal(t, expect, ref)
		}

		_, err = lkr.ResolveRef("he^^ad")
		require.Equal(t, err, ie.ErrNoSuchRef("he^^ad"))
	})
}
