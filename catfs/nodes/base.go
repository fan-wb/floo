package nodes

import (
	capnp_model "floo/catfs/nodes/capnp"
	"fmt"
	h "github.com/sahib/brig/util/hashlib"
	"time"
)

// Base is a place that holds all common attributes of all Nodes.
// It also defines some utility function that will be mixed into real nodes.
type Base struct {
	// Basename of this node
	name string

	// name of the user that last modified this node
	user string

	// Hash of this node (might be empty)
	tree h.Hash

	// Pointer hash to the content in the backend
	backend h.Hash

	// Content hash of this node
	content h.Hash

	// Last modification time of this node.
	modTime time.Time

	// Type of this node
	nodeType NodeType

	// Unique identifier for this node
	inode uint64
}

/////// UTILS /////////

func (b *Base) setBaseAttrsToNode(capnode capnp_model.Node) error {
	modTimeBin, err := b.modTime.MarshalBinary()
	if err != nil {
		return err
	}

	if err := capnode.SetModTime(string(modTimeBin)); err != nil {
		return err
	}
	if err := capnode.SetTreeHash(b.tree); err != nil {
		return err
	}
	if err := capnode.SetContentHash(b.content); err != nil {
		return err
	}
	if err := capnode.SetBackendHash(b.backend); err != nil {
		return err
	}
	if err := capnode.SetName(b.name); err != nil {
		return err
	}
	if err := capnode.SetUser(b.user); err != nil {
		return err
	}

	capnode.SetInode(b.inode)
	return nil
}

func (b *Base) parseBaseAttrsFromNode(capnode capnp_model.Node) error {
	var err error
	b.name, err = capnode.Name()
	if err != nil {
		return err
	}

	b.user, err = capnode.User()
	if err != nil {
		return err
	}

	b.tree, err = capnode.TreeHash()
	if err != nil {
		return err
	}

	b.content, err = capnode.ContentHash()
	if err != nil {
		return err
	}

	b.backend, err = capnode.BackendHash()
	if err != nil {
		return err
	}

	unparsedModTime, err := capnode.ModTime()
	if err != nil {
		return err
	}

	if err := b.modTime.UnmarshalBinary([]byte(unparsedModTime)); err != nil {
		return err
	}

	switch typ := capnode.Which(); typ {
	case capnp_model.Node_Which_file:
		b.nodeType = NodeTypeFile
	case capnp_model.Node_Which_directory:
		b.nodeType = NodeTypeDirectory
	case capnp_model.Node_Which_commit:
		b.nodeType = NodeTypeCommit
	case capnp_model.Node_Which_ghost:
		// Ghost set the nodeType themselves.
		// Ignore them here.
	default:
		return fmt.Errorf("bad capnp node type `%d`", typ)
	}

	b.inode = capnode.Inode()
	return nil
}
