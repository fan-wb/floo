package nodes

import (
	"capnproto.org/go/capnp/v3"
	capnp_model "floo/catfs/nodes/capnp"
	h "floo/util/hashlib"
	"fmt"
	"strings"
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

// copyBase will copy all attributes from the base.
func (b *Base) copyBase(inode uint64) Base {
	return Base{
		name:     b.name,
		user:     b.user,
		tree:     b.tree.Clone(),
		content:  b.content.Clone(),
		backend:  b.backend.Clone(),
		modTime:  b.modTime,
		nodeType: b.nodeType,
		inode:    inode,
	}
}

// User returns the user that last modified this node.
func (b *Base) User() string {
	return b.user
}

// Name returns the name of this node (e.g. /a/b/c -> c)
// The root directory will have the name empty string.
func (b *Base) Name() string {
	return b.name
}

// TreeHash returns the hash of this node.
func (b *Base) TreeHash() h.Hash {
	return b.tree
}

// ContentHash returns the content hash of this node.
func (b *Base) ContentHash() h.Hash {
	return b.content
}

// BackendHash returns the backend hash of this node.
func (b *Base) BackendHash() h.Hash {
	return b.backend
}

// Type returns the type of this node.
func (b *Base) Type() NodeType {
	return b.nodeType
}

// ModTime will return the last time this node's content
// was modified. Metadata changes are not recorded.
func (b *Base) ModTime() time.Time {
	return b.modTime
}

// Inode will return a unique ID that is different for each node.
func (b *Base) Inode() uint64 {
	return b.inode
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

func prefixSlash(s string) string {
	if !strings.HasPrefix(s, "/") {
		return "/" + s
	}

	return s
}

/////////////////////////////////////////
// MARSHAL HELPERS FOR ARBITRARY NODES //
/////////////////////////////////////////

// MarshalNode will convert any Node to a byte string
// Use UnmarshalNode to load a Node from it again.
func MarshalNode(nd Node) ([]byte, error) {
	msg, err := nd.ToCapnp()
	if err != nil {
		return nil, err
	}

	return msg.Marshal()
}

// UnmarshalNode will try to interpret data as a Node
func UnmarshalNode(data []byte) (Node, error) {
	msg, err := capnp.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	capNd, err := capnp_model.ReadRootNode(msg)
	if err != nil {
		return nil, err
	}

	return CapNodeToNode(capNd)
}

// CapNodeToNode converts a capnproto `capNd` to a normal `Node`.
func CapNodeToNode(capNd capnp_model.Node) (Node, error) {
	// Find out the correct node struct to initialize.
	var node Node

	switch typ := capNd.Which(); typ {
	case capnp_model.Node_Which_ghost:
		node = &Ghost{}
	case capnp_model.Node_Which_file:
		node = &File{}
	case capnp_model.Node_Which_directory:
		node = &Directory{}
	case capnp_model.Node_Which_commit:
		node = &Commit{}
	default:
		return nil, fmt.Errorf("Bad capnp node type `%d`", typ)
	}

	if err := node.FromCapnpNode(capNd); err != nil {
		return nil, err
	}

	return node, nil
}
