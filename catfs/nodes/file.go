package nodes

import (
	capnp "capnproto.org/go/capnp/v3"
	capnp_model "floo/catfs/nodes/capnp"
)

// File represents a single file in the repository.
// It stores all metadata about it and links to the actual data.
type File struct {
	Base

	size   uint64
	parent string
	key    []byte
}

// ToCapnp converts a file to a capnp message.
func (f *File) ToCapnp() (*capnp.Message, error) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, err
	}

	capNd, err := capnp_model.NewRootNode(seg)
	if err != nil {
		return nil, err
	}

	return msg, f.ToCapnpNode(seg, capNd)
}

// ToCapnpNode converts this node to a serializable capnp proto node.
func (f *File) ToCapnpNode(seg *capnp.Segment, capNd capnp_model.Node) error {
	if err := f.setBaseAttrsToNode(capNd); err != nil {
		return err
	}

	capFile, err := f.setFileAttrs(seg)
	if err != nil {
		return err
	}

	return capNd.SetFile(*capFile)
}

func (f *File) setFileAttrs(seg *capnp.Segment) (*capnp_model.File, error) {
	capFile, err := capnp_model.NewFile(seg)
	if err != nil {
		return nil, err
	}

	if err := capFile.SetParent(f.parent); err != nil {
		return nil, err
	}

	if err := capFile.SetKey(f.key); err != nil {
		return nil, err
	}

	capFile.SetSize(f.size)
	return &capFile, nil
}

// FromCapnpNode converts a serialized node to a normal node.
func (f *File) FromCapnpNode(capNd capnp_model.Node) error {
	if err := f.parseBaseAttrsFromNode(capNd); err != nil {
		return err
	}

	capFile, err := capNd.File()
	if err != nil {
		return err
	}

	return f.readFileAttrs(capFile)
}

func (f *File) readFileAttrs(capFile capnp_model.File) error {
	var err error

	f.parent, err = capFile.Parent()
	if err != nil {
		return err
	}

	f.nodeType = NodeTypeFile
	f.size = capFile.Size()
	f.key, err = capFile.Key()
	return err
}
