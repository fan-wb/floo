// Code generated by capnpc-go. DO NOT EDIT.

package capnp

import (
	capnp "capnproto.org/go/capnp/v3"
	text "capnproto.org/go/capnp/v3/encoding/text"
	schemas "capnproto.org/go/capnp/v3/schemas"
)

type User capnp.Struct

// User_TypeID is the unique identifier for the type User.
const User_TypeID = 0x861de4463c5a4a22

func NewUser(s *capnp.Segment) (User, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 5})
	return User(st), err
}

func NewRootUser(s *capnp.Segment) (User, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 5})
	return User(st), err
}

func ReadRootUser(msg *capnp.Message) (User, error) {
	root, err := msg.Root()
	return User(root.Struct()), err
}

func (s User) String() string {
	str, _ := text.Marshal(0x861de4463c5a4a22, capnp.Struct(s))
	return str
}

func (s User) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (User) DecodeFromPtr(p capnp.Ptr) User {
	return User(capnp.Struct{}.DecodeFromPtr(p))
}

func (s User) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s User) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s User) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s User) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s User) Name() (string, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.Text(), err
}

func (s User) HasName() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s User) NameBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.TextBytes(), err
}

func (s User) SetName(v string) error {
	return capnp.Struct(s).SetText(0, v)
}

func (s User) PasswordHash() (string, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return p.Text(), err
}

func (s User) HasPasswordHash() bool {
	return capnp.Struct(s).HasPtr(1)
}

func (s User) PasswordHashBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return p.TextBytes(), err
}

func (s User) SetPasswordHash(v string) error {
	return capnp.Struct(s).SetText(1, v)
}

func (s User) Salt() (string, error) {
	p, err := capnp.Struct(s).Ptr(2)
	return p.Text(), err
}

func (s User) HasSalt() bool {
	return capnp.Struct(s).HasPtr(2)
}

func (s User) SaltBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(2)
	return p.TextBytes(), err
}

func (s User) SetSalt(v string) error {
	return capnp.Struct(s).SetText(2, v)
}

func (s User) Folders() (capnp.TextList, error) {
	p, err := capnp.Struct(s).Ptr(3)
	return capnp.TextList(p.List()), err
}

func (s User) HasFolders() bool {
	return capnp.Struct(s).HasPtr(3)
}

func (s User) SetFolders(v capnp.TextList) error {
	return capnp.Struct(s).SetPtr(3, v.ToPtr())
}

// NewFolders sets the folders field to a newly
// allocated capnp.TextList, preferring placement in s's segment.
func (s User) NewFolders(n int32) (capnp.TextList, error) {
	l, err := capnp.NewTextList(capnp.Struct(s).Segment(), n)
	if err != nil {
		return capnp.TextList{}, err
	}
	err = capnp.Struct(s).SetPtr(3, l.ToPtr())
	return l, err
}

func (s User) Rights() (capnp.TextList, error) {
	p, err := capnp.Struct(s).Ptr(4)
	return capnp.TextList(p.List()), err
}

func (s User) HasRights() bool {
	return capnp.Struct(s).HasPtr(4)
}

func (s User) SetRights(v capnp.TextList) error {
	return capnp.Struct(s).SetPtr(4, v.ToPtr())
}

// NewRights sets the rights field to a newly
// allocated capnp.TextList, preferring placement in s's segment.
func (s User) NewRights(n int32) (capnp.TextList, error) {
	l, err := capnp.NewTextList(capnp.Struct(s).Segment(), n)
	if err != nil {
		return capnp.TextList{}, err
	}
	err = capnp.Struct(s).SetPtr(4, l.ToPtr())
	return l, err
}

// User_List is a list of User.
type User_List = capnp.StructList[User]

// NewUser creates a new list of User.
func NewUser_List(s *capnp.Segment, sz int32) (User_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 5}, sz)
	return capnp.StructList[User](l), err
}

// User_Future is a wrapper for a User promised by a client call.
type User_Future struct{ *capnp.Future }

func (p User_Future) Struct() (User, error) {
	s, err := p.Future.Struct()
	return User(s), err
}

const schema_a0b1c18bd0f965c4 = "x\xda\\\xca\xb1J3A\x00\xc4\xf1\x99\xdd\xcd\xf7\x81" +
	"\xc4\x9c\x0b[*\x82\xa5\xa0G\xda (\x16\"V\xd9" +
	"\xc2\xc6n\xf5\xd6D\x89\xc9q{\x12-D\x85 \x8a" +
	"\x0a\x96\x16\x16\x0a\xbe\x80\x9d\x9d\x08\xdak\xe1\x1b\xf8\x12" +
	"V'\x17Hc\xf7\x9f\x1f3q\xbb\xa4\xea\xe3/\x84" +
	"\xb0\xa6\xf2\xaf\x98Y\xdbXX\xf9\x9e:\x83\x9ed\xf1" +
	"\xe6\x7f>._\x9f\xeeQ\xa9\xfc\x07\xea\xefc\xd4_" +
	"e|N\x13sE\xcb\xe5\xbe\xef\x0ec\x99l\xc6[" +
	".\xed\xa6\xf1~\xf0\xd9\xfc0\x1b\xeb\xc1g@\x93\xb4" +
	"F*@\x11\xd0G\xb3\x80=\x90\xb4\x03AM\x1a\x96" +
	"x\xba\x0b\xd8\x13I{%\xa8\x850\x14\x80\xbe(\x9f" +
	"\x03I{#\xa8\xa54\x94\x80\xbe^\x06\xec\xb9\xa4}" +
	"\x14\xd4J\x19*@?4\x00{'i\x9f\x05\xa3\xae" +
	"\xdb\xf3\xacB\xb0\x0a\x16\xa9\x0b\xa1\xdf\xcb\x12D\xab." +
	"\xb4G\x1c\x05\xd7\xc9G\xe3x\xbb\xd7I|\x16X\x03" +
	"\x9b\x92C\xae\x81\x8b\xd9N\xab\x9d\xff\xd5\xdf\x00\x00\x00" +
	"\xff\xff\x99\xed=\x1f"

func init() {
	schemas.Register(schema_a0b1c18bd0f965c4,
		0x861de4463c5a4a22)
}
