package nodes

import (
	h "floo/util/hashlib"
)

// Linker will tell a node how it relates to other nodes
// and gives it the ability to resolve other nodes by hash.
// Apart from that it gives the underlying linker implementation
// the possibility to be notified when a hash changes.
type Linker interface {
	// Root should return the current root directory.
	Root() (*Directory, error)

	// LookupNode should resolve `path` starting from the root directory.
	// If the path does not exist an error is returned and can be checked
	// with IsNoSuchFileError()
	LookupNode(path string) (Node, error)

	// NodeByHash resolves the hash to a specific node.
	// If the node does not exist, nil is returned.
	NodeByHash(hash h.Hash) (Node, error)

	// MemIndexSwap should be called when
	// the hash of a node changes.
	MemIndexSwap(nd Node, oldHash h.Hash, updatePathIndex bool)

	// MemSetRoot should be called when the current root directory changed.
	MemSetRoot(root *Directory)
}
