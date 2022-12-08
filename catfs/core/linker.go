package core

import (
	"floo/catfs/db"
	ie "floo/catfs/errors"
	n "floo/catfs/nodes"
	h "floo/util/hashlib"
	"floo/util/trie"
	"strings"
)

// Linker implements the basic logic of floo's data model
// It uses an underlying key/value database to
// store a Merkle-DAG with versioned metadata,
// similar to what git does internally.
type Linker struct {
	kv db.Database

	// root of the filesystem
	root *n.Directory

	// Path lookup trie
	ptrie *trie.Node

	// B58Hash to node
	index map[string]n.Node

	// UID to node
	inodeIndex map[uint64]n.Node

	// Cache for the linker owner.
	owner string
}

// NewLinker returns a new lkr, ready to use. It assumes the key value store
// is working and does no check on this.
func NewLinker(kv db.Database) *Linker {
	lkr := &Linker{kv: kv}
	lkr.MemIndexClear()
	return lkr
}

// MemIndexAdd adds `nd` to the in memory index.
func (lkr *Linker) MemIndexAdd(nd n.Node, updatePathIndex bool) {
	lkr.index[nd.TreeHash().B58String()] = nd
	lkr.inodeIndex[nd.Inode()] = nd

	if updatePathIndex {
		path := nd.Path()
		if nd.Type() == n.NodeTypeDirectory {
			path = appendDot(path)
		}
		lkr.ptrie.InsertWithData(path, nd)
	}
}

// MemIndexSwap updates an entry of the in memory index, by deleting
// the old entry referenced by oldHash (maybe nil). This is necessary
// to ensure that old hashes do not resolve to the new, updated instance.
// If the old instance is needed, it will be loaded as new instance.
// You should not need to call this function, except when implementing own Nodes.
func (lkr *Linker) MemIndexSwap(nd n.Node, oldHash h.Hash, updatePathIndex bool) {
	if oldHash != nil {
		delete(lkr.index, oldHash.B58String())
	}

	lkr.MemIndexAdd(nd, updatePathIndex)
}

// MemSetRoot sets the current root, but does not store it yet. It's supposed
// to be called after in-memory modifications. Only implementors of new Nodes
// might need to call this function.
func (lkr *Linker) MemSetRoot(root *n.Directory) {
	if lkr.root != nil {
		lkr.MemIndexSwap(root, lkr.root.TreeHash(), true)
	} else {
		lkr.MemIndexAdd(root, true)
	}

	lkr.root = root
}

// MemIndexPurge removes `nd` from the memory index.
func (lkr *Linker) MemIndexPurge(nd n.Node) {
	delete(lkr.inodeIndex, nd.Inode())
	delete(lkr.index, nd.TreeHash().B58String())
	lkr.ptrie.Lookup(nd.Path()).Remove()
}

// MemIndexClear resets the memory index to zero.
// This should not be called mid-flight in operations,
// but should be okay to call between atomic operations.
func (lkr *Linker) MemIndexClear() {
	lkr.ptrie = trie.NewNode()
	lkr.index = make(map[string]n.Node)
	lkr.inodeIndex = make(map[uint64]n.Node)
	lkr.root = nil
}

//////////////////////////
// COMMON NODE HANDLING //
//////////////////////////

// NextInode returns a unique identifier, used to identify a single node. You
// should not need to call this function, except when implementing own nodes.
func (lkr *Linker) NextInode() uint64 {
	// TODO
	return uint64(0)
}

// FilesByContents checks what files are associated with the content hashes in
// `contents`. It returns a map of content hash b58 to file. This method is
// quite heavy and should not be used in loops. There is room for optimizations.
func (lkr *Linker) FilesByContents(contents []h.Hash) (map[string]*n.File, error) {
	keys, err := lkr.kv.Keys()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*n.File)
	for _, key := range keys {
		// Filter non-node storage:
		fullKey := strings.Join(key, "/")
		if !strings.HasPrefix(fullKey, "objects") &&
			!strings.HasPrefix(fullKey, "stage/objects") {
			continue
		}

		data, err := lkr.kv.Get(key...)
		if err != nil {
			return nil, err
		}

		nd, err := n.UnmarshalNode(data)
		if err != nil {
			return nil, err
		}

		if nd.Type() != n.NodeTypeFile {
			continue
		}

		file, ok := nd.(*n.File)
		if !ok {
			return nil, ie.ErrBadNode
		}

		for _, content := range contents {
			if content.Equal(file.BackendHash()) {
				result[content.B58String()] = file
			}
		}
	}

	return result, nil
}

func appendDot(path string) string {
	// path.Join() calls path.Clean() which in turn
	// removes the '.' at the end when trying to join that.
	// But since we use the dot to mark directories we shouldn't do that.
	if strings.HasSuffix(path, "/") {
		return path + "."
	}

	return path + "/."
}
