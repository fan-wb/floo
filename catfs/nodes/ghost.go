package nodes

import (
	ie "floo/catfs/errors"
	h "floo/util/hashlib"
	"fmt"
)

// Ghost is a special kind of Node that marks a moved node.
// If a file was moved, a ghost will be created for the old place.
// If another file is moved to the new place, the ghost will be "resurrected"
// with the new content.
type Ghost struct {
	ModNode

	ghostPath  string
	ghostInode uint64
	oldType    NodeType
}

// Type always returns NodeTypeGhost
func (g *Ghost) Type() NodeType {
	return NodeTypeGhost
}

// OldNode returns the node the ghost was when it still was alive.
func (g *Ghost) OldNode() ModNode {
	return g.ModNode
}

// OldFile returns the file the ghost was when it still was alive.
// Returns ErrBadNode when it wasn't a file.
func (g *Ghost) OldFile() (*File, error) {
	file, ok := g.ModNode.(*File)
	if !ok {
		return nil, ie.ErrBadNode
	}

	return file, nil
}

// OldDirectory returns the old directory that the node was in lifetime
// If the ghost was not a directory, ErrBadNode is returned.
func (g *Ghost) OldDirectory() (*Directory, error) {
	directory, ok := g.ModNode.(*Directory)
	if !ok {
		return nil, ie.ErrBadNode
	}

	return directory, nil
}

func (g *Ghost) String() string {
	return fmt.Sprintf("<ghost: %s %v>", g.TreeHash(), g.ModNode)
}

// Path returns the path of the node.
func (g *Ghost) Path() string {
	return g.ghostPath
}

// TreeHash returns the hash of the node.
func (g *Ghost) TreeHash() h.Hash {
	return h.Sum([]byte(fmt.Sprintf("ghost:%s", g.ModNode.TreeHash())))
}

// Inode returns the inode
func (g *Ghost) Inode() uint64 {
	return g.ghostInode
}

// SetGhostPath sets the path of the ghost.
func (g *Ghost) SetGhostPath(newPath string) {
	g.ghostPath = newPath
}
