// Code generated by capnpc-go. DO NOT EDIT.

package capnp

import (
	capnp "capnproto.org/go/capnp/v3"
	text "capnproto.org/go/capnp/v3/encoding/text"
	schemas "capnproto.org/go/capnp/v3/schemas"
)

type Hint capnp.Struct

// Hint_TypeID is the unique identifier for the type Hint.
const Hint_TypeID = 0xb2ec3fe21ddc803f

func NewHint(s *capnp.Segment) (Hint, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 3})
	return Hint(st), err
}

func NewRootHint(s *capnp.Segment) (Hint, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 3})
	return Hint(st), err
}

func ReadRootHint(msg *capnp.Message) (Hint, error) {
	root, err := msg.Root()
	return Hint(root.Struct()), err
}

func (s Hint) String() string {
	str, _ := text.Marshal(0xb2ec3fe21ddc803f, capnp.Struct(s))
	return str
}

func (s Hint) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (Hint) DecodeFromPtr(p capnp.Ptr) Hint {
	return Hint(capnp.Struct{}.DecodeFromPtr(p))
}

func (s Hint) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s Hint) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s Hint) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s Hint) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s Hint) Path() (string, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.Text(), err
}

func (s Hint) HasPath() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s Hint) PathBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.TextBytes(), err
}

func (s Hint) SetPath(v string) error {
	return capnp.Struct(s).SetText(0, v)
}

func (s Hint) EncryptionAlgo() (string, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return p.Text(), err
}

func (s Hint) HasEncryptionAlgo() bool {
	return capnp.Struct(s).HasPtr(1)
}

func (s Hint) EncryptionAlgoBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return p.TextBytes(), err
}

func (s Hint) SetEncryptionAlgo(v string) error {
	return capnp.Struct(s).SetText(1, v)
}

func (s Hint) CompressionAlgo() (string, error) {
	p, err := capnp.Struct(s).Ptr(2)
	return p.Text(), err
}

func (s Hint) HasCompressionAlgo() bool {
	return capnp.Struct(s).HasPtr(2)
}

func (s Hint) CompressionAlgoBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(2)
	return p.TextBytes(), err
}

func (s Hint) SetCompressionAlgo(v string) error {
	return capnp.Struct(s).SetText(2, v)
}

// Hint_List is a list of Hint.
type Hint_List = capnp.StructList[Hint]

// NewHint creates a new list of Hint.
func NewHint_List(s *capnp.Segment, sz int32) (Hint_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 3}, sz)
	return capnp.StructList[Hint](l), err
}

// Hint_Future is a wrapper for a Hint promised by a client call.
type Hint_Future struct{ *capnp.Future }

func (p Hint_Future) Struct() (Hint, error) {
	s, err := p.Future.Struct()
	return Hint(s), err
}

const schema_ea883e7d5248d81b = "x\xda<\xc91J\xc4@\x18\x86\xe1\xef\xfb'\xeb6" +
	"\x11\x1dI'{\x01A\x97m-\xccZ\x08\x11,\xf6" +
	"\x17KA\xc2\x104\x10\x93!\x13\x04\x05A\xac\xbc\x8d" +
	"\xe0\x05\xac<\x80\x85`ae)\xdead\x11-\xdf" +
	"\xf7Y\x0f\xf3d\xb6\xfaL\x88f\xa3\x95\x98\xdf}L" +
	">\xf3\xef'\xd8\x09\xe3\xe6{q|\xbb\xf7\xf0\x85\x91" +
	"\x19\x03\xb3\x97\x0d\xda\xb71`_\x1f\xb1\x1dC\xd5_" +
	"U\xfd\xd4\x99\xd2\xb7~\xdat\xael\xceJ_\xef\xb8" +
	"e\xef\x16u\xcbaAjj\x12 !`\x0f\xb6\x00" +
	"\x9d\x1b\xea\x91\xd0\x92\x19\x97\xf3\xf0\x06\xd0\xc2PO\x84" +
	"V$\xa3\x00V\xef\x01]\x18\xea\xa9p\xcd\x97\xc3\x05" +
	"S\x08S0V\xad\xeb\xaf\xfdP#\xef\xda\xfd\xe6\xbc" +
	"\xfb\x07\xd7]\xfa\xbe\x0a\x81\xf5/\xe0O~\x02\x00\x00" +
	"\xff\xff\x8bO5\x7f"

func init() {
	schemas.Register(schema_ea883e7d5248d81b,
		0xb2ec3fe21ddc803f)
}
