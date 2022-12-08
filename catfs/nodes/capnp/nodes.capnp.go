// Code generated by capnpc-go. DO NOT EDIT.

package capnp

import (
	capnp "capnproto.org/go/capnp/v3"
	text "capnproto.org/go/capnp/v3/encoding/text"
	schemas "capnproto.org/go/capnp/v3/schemas"
	strconv "strconv"
)

// Commit is a set of changes to nodes
type Commit capnp.Struct
type Commit_merge Commit

// Commit_TypeID is the unique identifier for the type Commit.
const Commit_TypeID = 0x8da013c66e545daf

func NewCommit(s *capnp.Segment) (Commit, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 6})
	return Commit(st), err
}

func NewRootCommit(s *capnp.Segment) (Commit, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 6})
	return Commit(st), err
}

func ReadRootCommit(msg *capnp.Message) (Commit, error) {
	root, err := msg.Root()
	return Commit(root.Struct()), err
}

func (s Commit) String() string {
	str, _ := text.Marshal(0x8da013c66e545daf, capnp.Struct(s))
	return str
}

func (s Commit) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (Commit) DecodeFromPtr(p capnp.Ptr) Commit {
	return Commit(capnp.Struct{}.DecodeFromPtr(p))
}

func (s Commit) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s Commit) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s Commit) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s Commit) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s Commit) Msg() (string, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.Text(), err
}

func (s Commit) HasMsg() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s Commit) MsgBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.TextBytes(), err
}

func (s Commit) SetMsg(v string) error {
	return capnp.Struct(s).SetText(0, v)
}

func (s Commit) Author() (string, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return p.Text(), err
}

func (s Commit) HasAuthor() bool {
	return capnp.Struct(s).HasPtr(1)
}

func (s Commit) AuthorBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return p.TextBytes(), err
}

func (s Commit) SetAuthor(v string) error {
	return capnp.Struct(s).SetText(1, v)
}

func (s Commit) Parent() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(2)
	return []byte(p.Data()), err
}

func (s Commit) HasParent() bool {
	return capnp.Struct(s).HasPtr(2)
}

func (s Commit) SetParent(v []byte) error {
	return capnp.Struct(s).SetData(2, v)
}

func (s Commit) Root() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(3)
	return []byte(p.Data()), err
}

func (s Commit) HasRoot() bool {
	return capnp.Struct(s).HasPtr(3)
}

func (s Commit) SetRoot(v []byte) error {
	return capnp.Struct(s).SetData(3, v)
}

func (s Commit) Index() int64 {
	return int64(capnp.Struct(s).Uint64(0))
}

func (s Commit) SetIndex(v int64) {
	capnp.Struct(s).SetUint64(0, uint64(v))
}

func (s Commit) Merge() Commit_merge { return Commit_merge(s) }

func (s Commit_merge) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s Commit_merge) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s Commit_merge) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s Commit_merge) With() (string, error) {
	p, err := capnp.Struct(s).Ptr(4)
	return p.Text(), err
}

func (s Commit_merge) HasWith() bool {
	return capnp.Struct(s).HasPtr(4)
}

func (s Commit_merge) WithBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(4)
	return p.TextBytes(), err
}

func (s Commit_merge) SetWith(v string) error {
	return capnp.Struct(s).SetText(4, v)
}

func (s Commit_merge) Head() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(5)
	return []byte(p.Data()), err
}

func (s Commit_merge) HasHead() bool {
	return capnp.Struct(s).HasPtr(5)
}

func (s Commit_merge) SetHead(v []byte) error {
	return capnp.Struct(s).SetData(5, v)
}

// Commit_List is a list of Commit.
type Commit_List = capnp.StructList[Commit]

// NewCommit creates a new list of Commit.
func NewCommit_List(s *capnp.Segment, sz int32) (Commit_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 6}, sz)
	return capnp.StructList[Commit](l), err
}

// Commit_Future is a wrapper for a Commit promised by a client call.
type Commit_Future struct{ *capnp.Future }

func (f Commit_Future) Struct() (Commit, error) {
	p, err := f.Future.Ptr()
	return Commit(p.Struct()), err
}
func (p Commit_Future) Merge() Commit_merge_Future { return Commit_merge_Future{p.Future} }

// Commit_merge_Future is a wrapper for a Commit_merge promised by a client call.
type Commit_merge_Future struct{ *capnp.Future }

func (f Commit_merge_Future) Struct() (Commit_merge, error) {
	p, err := f.Future.Ptr()
	return Commit_merge(p.Struct()), err
}

// A single directory entry
type DirEntry capnp.Struct

// DirEntry_TypeID is the unique identifier for the type DirEntry.
const DirEntry_TypeID = 0x8b15ee76774b1f9d

func NewDirEntry(s *capnp.Segment) (DirEntry, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return DirEntry(st), err
}

func NewRootDirEntry(s *capnp.Segment) (DirEntry, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return DirEntry(st), err
}

func ReadRootDirEntry(msg *capnp.Message) (DirEntry, error) {
	root, err := msg.Root()
	return DirEntry(root.Struct()), err
}

func (s DirEntry) String() string {
	str, _ := text.Marshal(0x8b15ee76774b1f9d, capnp.Struct(s))
	return str
}

func (s DirEntry) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (DirEntry) DecodeFromPtr(p capnp.Ptr) DirEntry {
	return DirEntry(capnp.Struct{}.DecodeFromPtr(p))
}

func (s DirEntry) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s DirEntry) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s DirEntry) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s DirEntry) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s DirEntry) Name() (string, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.Text(), err
}

func (s DirEntry) HasName() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s DirEntry) NameBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.TextBytes(), err
}

func (s DirEntry) SetName(v string) error {
	return capnp.Struct(s).SetText(0, v)
}

func (s DirEntry) Hash() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return []byte(p.Data()), err
}

func (s DirEntry) HasHash() bool {
	return capnp.Struct(s).HasPtr(1)
}

func (s DirEntry) SetHash(v []byte) error {
	return capnp.Struct(s).SetData(1, v)
}

// DirEntry_List is a list of DirEntry.
type DirEntry_List = capnp.StructList[DirEntry]

// NewDirEntry creates a new list of DirEntry.
func NewDirEntry_List(s *capnp.Segment, sz int32) (DirEntry_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	return capnp.StructList[DirEntry](l), err
}

// DirEntry_Future is a wrapper for a DirEntry promised by a client call.
type DirEntry_Future struct{ *capnp.Future }

func (f DirEntry_Future) Struct() (DirEntry, error) {
	p, err := f.Future.Ptr()
	return DirEntry(p.Struct()), err
}

// Directory contains one or more directories or files
type Directory capnp.Struct

// Directory_TypeID is the unique identifier for the type Directory.
const Directory_TypeID = 0xe24c59306c829c01

func NewDirectory(s *capnp.Segment) (Directory, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 3})
	return Directory(st), err
}

func NewRootDirectory(s *capnp.Segment) (Directory, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 3})
	return Directory(st), err
}

func ReadRootDirectory(msg *capnp.Message) (Directory, error) {
	root, err := msg.Root()
	return Directory(root.Struct()), err
}

func (s Directory) String() string {
	str, _ := text.Marshal(0xe24c59306c829c01, capnp.Struct(s))
	return str
}

func (s Directory) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (Directory) DecodeFromPtr(p capnp.Ptr) Directory {
	return Directory(capnp.Struct{}.DecodeFromPtr(p))
}

func (s Directory) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s Directory) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s Directory) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s Directory) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s Directory) Size() uint64 {
	return capnp.Struct(s).Uint64(0)
}

func (s Directory) SetSize(v uint64) {
	capnp.Struct(s).SetUint64(0, v)
}

func (s Directory) Parent() (string, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.Text(), err
}

func (s Directory) HasParent() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s Directory) ParentBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.TextBytes(), err
}

func (s Directory) SetParent(v string) error {
	return capnp.Struct(s).SetText(0, v)
}

func (s Directory) Children() (DirEntry_List, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return DirEntry_List(p.List()), err
}

func (s Directory) HasChildren() bool {
	return capnp.Struct(s).HasPtr(1)
}

func (s Directory) SetChildren(v DirEntry_List) error {
	return capnp.Struct(s).SetPtr(1, v.ToPtr())
}

// NewChildren sets the children field to a newly
// allocated DirEntry_List, preferring placement in s's segment.
func (s Directory) NewChildren(n int32) (DirEntry_List, error) {
	l, err := NewDirEntry_List(capnp.Struct(s).Segment(), n)
	if err != nil {
		return DirEntry_List{}, err
	}
	err = capnp.Struct(s).SetPtr(1, l.ToPtr())
	return l, err
}
func (s Directory) Contents() (DirEntry_List, error) {
	p, err := capnp.Struct(s).Ptr(2)
	return DirEntry_List(p.List()), err
}

func (s Directory) HasContents() bool {
	return capnp.Struct(s).HasPtr(2)
}

func (s Directory) SetContents(v DirEntry_List) error {
	return capnp.Struct(s).SetPtr(2, v.ToPtr())
}

// NewContents sets the contents field to a newly
// allocated DirEntry_List, preferring placement in s's segment.
func (s Directory) NewContents(n int32) (DirEntry_List, error) {
	l, err := NewDirEntry_List(capnp.Struct(s).Segment(), n)
	if err != nil {
		return DirEntry_List{}, err
	}
	err = capnp.Struct(s).SetPtr(2, l.ToPtr())
	return l, err
}

// Directory_List is a list of Directory.
type Directory_List = capnp.StructList[Directory]

// NewDirectory creates a new list of Directory.
func NewDirectory_List(s *capnp.Segment, sz int32) (Directory_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 3}, sz)
	return capnp.StructList[Directory](l), err
}

// Directory_Future is a wrapper for a Directory promised by a client call.
type Directory_Future struct{ *capnp.Future }

func (f Directory_Future) Struct() (Directory, error) {
	p, err := f.Future.Ptr()
	return Directory(p.Struct()), err
}

// A leaf node in the MDAG
type File capnp.Struct

// File_TypeID is the unique identifier for the type File.
const File_TypeID = 0x8ea7393d37893155

func NewFile(s *capnp.Segment) (File, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	return File(st), err
}

func NewRootFile(s *capnp.Segment) (File, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	return File(st), err
}

func ReadRootFile(msg *capnp.Message) (File, error) {
	root, err := msg.Root()
	return File(root.Struct()), err
}

func (s File) String() string {
	str, _ := text.Marshal(0x8ea7393d37893155, capnp.Struct(s))
	return str
}

func (s File) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (File) DecodeFromPtr(p capnp.Ptr) File {
	return File(capnp.Struct{}.DecodeFromPtr(p))
}

func (s File) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s File) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s File) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s File) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s File) Size() uint64 {
	return capnp.Struct(s).Uint64(0)
}

func (s File) SetSize(v uint64) {
	capnp.Struct(s).SetUint64(0, v)
}

func (s File) Parent() (string, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.Text(), err
}

func (s File) HasParent() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s File) ParentBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.TextBytes(), err
}

func (s File) SetParent(v string) error {
	return capnp.Struct(s).SetText(0, v)
}

func (s File) Key() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return []byte(p.Data()), err
}

func (s File) HasKey() bool {
	return capnp.Struct(s).HasPtr(1)
}

func (s File) SetKey(v []byte) error {
	return capnp.Struct(s).SetData(1, v)
}

// File_List is a list of File.
type File_List = capnp.StructList[File]

// NewFile creates a new list of File.
func NewFile_List(s *capnp.Segment, sz int32) (File_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2}, sz)
	return capnp.StructList[File](l), err
}

// File_Future is a wrapper for a File promised by a client call.
type File_Future struct{ *capnp.Future }

func (f File_Future) Struct() (File, error) {
	p, err := f.Future.Ptr()
	return File(p.Struct()), err
}

// Ghost indicates that a certain node was at this path once
type Ghost capnp.Struct
type Ghost_Which uint16

const (
	Ghost_Which_commit    Ghost_Which = 0
	Ghost_Which_directory Ghost_Which = 1
	Ghost_Which_file      Ghost_Which = 2
)

func (w Ghost_Which) String() string {
	const s = "commitdirectoryfile"
	switch w {
	case Ghost_Which_commit:
		return s[0:6]
	case Ghost_Which_directory:
		return s[6:15]
	case Ghost_Which_file:
		return s[15:19]

	}
	return "Ghost_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

// Ghost_TypeID is the unique identifier for the type Ghost.
const Ghost_TypeID = 0x80c828d7e89c12ea

func NewGhost(s *capnp.Segment) (Ghost, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 2})
	return Ghost(st), err
}

func NewRootGhost(s *capnp.Segment) (Ghost, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 2})
	return Ghost(st), err
}

func ReadRootGhost(msg *capnp.Message) (Ghost, error) {
	root, err := msg.Root()
	return Ghost(root.Struct()), err
}

func (s Ghost) String() string {
	str, _ := text.Marshal(0x80c828d7e89c12ea, capnp.Struct(s))
	return str
}

func (s Ghost) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (Ghost) DecodeFromPtr(p capnp.Ptr) Ghost {
	return Ghost(capnp.Struct{}.DecodeFromPtr(p))
}

func (s Ghost) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}

func (s Ghost) Which() Ghost_Which {
	return Ghost_Which(capnp.Struct(s).Uint16(8))
}
func (s Ghost) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s Ghost) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s Ghost) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s Ghost) GhostInode() uint64 {
	return capnp.Struct(s).Uint64(0)
}

func (s Ghost) SetGhostInode(v uint64) {
	capnp.Struct(s).SetUint64(0, v)
}

func (s Ghost) GhostPath() (string, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.Text(), err
}

func (s Ghost) HasGhostPath() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s Ghost) GhostPathBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.TextBytes(), err
}

func (s Ghost) SetGhostPath(v string) error {
	return capnp.Struct(s).SetText(0, v)
}

func (s Ghost) Commit() (Commit, error) {
	if capnp.Struct(s).Uint16(8) != 0 {
		panic("Which() != commit")
	}
	p, err := capnp.Struct(s).Ptr(1)
	return Commit(p.Struct()), err
}

func (s Ghost) HasCommit() bool {
	if capnp.Struct(s).Uint16(8) != 0 {
		return false
	}
	return capnp.Struct(s).HasPtr(1)
}

func (s Ghost) SetCommit(v Commit) error {
	capnp.Struct(s).SetUint16(8, 0)
	return capnp.Struct(s).SetPtr(1, capnp.Struct(v).ToPtr())
}

// NewCommit sets the commit field to a newly
// allocated Commit struct, preferring placement in s's segment.
func (s Ghost) NewCommit() (Commit, error) {
	capnp.Struct(s).SetUint16(8, 0)
	ss, err := NewCommit(capnp.Struct(s).Segment())
	if err != nil {
		return Commit{}, err
	}
	err = capnp.Struct(s).SetPtr(1, capnp.Struct(ss).ToPtr())
	return ss, err
}

func (s Ghost) Directory() (Directory, error) {
	if capnp.Struct(s).Uint16(8) != 1 {
		panic("Which() != directory")
	}
	p, err := capnp.Struct(s).Ptr(1)
	return Directory(p.Struct()), err
}

func (s Ghost) HasDirectory() bool {
	if capnp.Struct(s).Uint16(8) != 1 {
		return false
	}
	return capnp.Struct(s).HasPtr(1)
}

func (s Ghost) SetDirectory(v Directory) error {
	capnp.Struct(s).SetUint16(8, 1)
	return capnp.Struct(s).SetPtr(1, capnp.Struct(v).ToPtr())
}

// NewDirectory sets the directory field to a newly
// allocated Directory struct, preferring placement in s's segment.
func (s Ghost) NewDirectory() (Directory, error) {
	capnp.Struct(s).SetUint16(8, 1)
	ss, err := NewDirectory(capnp.Struct(s).Segment())
	if err != nil {
		return Directory{}, err
	}
	err = capnp.Struct(s).SetPtr(1, capnp.Struct(ss).ToPtr())
	return ss, err
}

func (s Ghost) File() (File, error) {
	if capnp.Struct(s).Uint16(8) != 2 {
		panic("Which() != file")
	}
	p, err := capnp.Struct(s).Ptr(1)
	return File(p.Struct()), err
}

func (s Ghost) HasFile() bool {
	if capnp.Struct(s).Uint16(8) != 2 {
		return false
	}
	return capnp.Struct(s).HasPtr(1)
}

func (s Ghost) SetFile(v File) error {
	capnp.Struct(s).SetUint16(8, 2)
	return capnp.Struct(s).SetPtr(1, capnp.Struct(v).ToPtr())
}

// NewFile sets the file field to a newly
// allocated File struct, preferring placement in s's segment.
func (s Ghost) NewFile() (File, error) {
	capnp.Struct(s).SetUint16(8, 2)
	ss, err := NewFile(capnp.Struct(s).Segment())
	if err != nil {
		return File{}, err
	}
	err = capnp.Struct(s).SetPtr(1, capnp.Struct(ss).ToPtr())
	return ss, err
}

// Ghost_List is a list of Ghost.
type Ghost_List = capnp.StructList[Ghost]

// NewGhost creates a new list of Ghost.
func NewGhost_List(s *capnp.Segment, sz int32) (Ghost_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 2}, sz)
	return capnp.StructList[Ghost](l), err
}

// Ghost_Future is a wrapper for a Ghost promised by a client call.
type Ghost_Future struct{ *capnp.Future }

func (f Ghost_Future) Struct() (Ghost, error) {
	p, err := f.Future.Ptr()
	return Ghost(p.Struct()), err
}
func (p Ghost_Future) Commit() Commit_Future {
	return Commit_Future{Future: p.Future.Field(1, nil)}
}
func (p Ghost_Future) Directory() Directory_Future {
	return Directory_Future{Future: p.Future.Field(1, nil)}
}
func (p Ghost_Future) File() File_Future {
	return File_Future{Future: p.Future.Field(1, nil)}
}

// Node is a node in the merkle dag of floo
type Node capnp.Struct
type Node_Which uint16

const (
	Node_Which_commit    Node_Which = 0
	Node_Which_directory Node_Which = 1
	Node_Which_file      Node_Which = 2
	Node_Which_ghost     Node_Which = 3
)

func (w Node_Which) String() string {
	const s = "commitdirectoryfileghost"
	switch w {
	case Node_Which_commit:
		return s[0:6]
	case Node_Which_directory:
		return s[6:15]
	case Node_Which_file:
		return s[15:19]
	case Node_Which_ghost:
		return s[19:24]

	}
	return "Node_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

// Node_TypeID is the unique identifier for the type Node.
const Node_TypeID = 0xa629eb7f7066fae3

func NewNode(s *capnp.Segment) (Node, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 7})
	return Node(st), err
}

func NewRootNode(s *capnp.Segment) (Node, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 7})
	return Node(st), err
}

func ReadRootNode(msg *capnp.Message) (Node, error) {
	root, err := msg.Root()
	return Node(root.Struct()), err
}

func (s Node) String() string {
	str, _ := text.Marshal(0xa629eb7f7066fae3, capnp.Struct(s))
	return str
}

func (s Node) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (Node) DecodeFromPtr(p capnp.Ptr) Node {
	return Node(capnp.Struct{}.DecodeFromPtr(p))
}

func (s Node) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}

func (s Node) Which() Node_Which {
	return Node_Which(capnp.Struct(s).Uint16(8))
}
func (s Node) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s Node) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s Node) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s Node) Name() (string, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.Text(), err
}

func (s Node) HasName() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s Node) NameBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.TextBytes(), err
}

func (s Node) SetName(v string) error {
	return capnp.Struct(s).SetText(0, v)
}

func (s Node) TreeHash() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return []byte(p.Data()), err
}

func (s Node) HasTreeHash() bool {
	return capnp.Struct(s).HasPtr(1)
}

func (s Node) SetTreeHash(v []byte) error {
	return capnp.Struct(s).SetData(1, v)
}

func (s Node) ModTime() (string, error) {
	p, err := capnp.Struct(s).Ptr(2)
	return p.Text(), err
}

func (s Node) HasModTime() bool {
	return capnp.Struct(s).HasPtr(2)
}

func (s Node) ModTimeBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(2)
	return p.TextBytes(), err
}

func (s Node) SetModTime(v string) error {
	return capnp.Struct(s).SetText(2, v)
}

func (s Node) Inode() uint64 {
	return capnp.Struct(s).Uint64(0)
}

func (s Node) SetInode(v uint64) {
	capnp.Struct(s).SetUint64(0, v)
}

func (s Node) ContentHash() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(3)
	return []byte(p.Data()), err
}

func (s Node) HasContentHash() bool {
	return capnp.Struct(s).HasPtr(3)
}

func (s Node) SetContentHash(v []byte) error {
	return capnp.Struct(s).SetData(3, v)
}

func (s Node) User() (string, error) {
	p, err := capnp.Struct(s).Ptr(4)
	return p.Text(), err
}

func (s Node) HasUser() bool {
	return capnp.Struct(s).HasPtr(4)
}

func (s Node) UserBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(4)
	return p.TextBytes(), err
}

func (s Node) SetUser(v string) error {
	return capnp.Struct(s).SetText(4, v)
}

func (s Node) Commit() (Commit, error) {
	if capnp.Struct(s).Uint16(8) != 0 {
		panic("Which() != commit")
	}
	p, err := capnp.Struct(s).Ptr(5)
	return Commit(p.Struct()), err
}

func (s Node) HasCommit() bool {
	if capnp.Struct(s).Uint16(8) != 0 {
		return false
	}
	return capnp.Struct(s).HasPtr(5)
}

func (s Node) SetCommit(v Commit) error {
	capnp.Struct(s).SetUint16(8, 0)
	return capnp.Struct(s).SetPtr(5, capnp.Struct(v).ToPtr())
}

// NewCommit sets the commit field to a newly
// allocated Commit struct, preferring placement in s's segment.
func (s Node) NewCommit() (Commit, error) {
	capnp.Struct(s).SetUint16(8, 0)
	ss, err := NewCommit(capnp.Struct(s).Segment())
	if err != nil {
		return Commit{}, err
	}
	err = capnp.Struct(s).SetPtr(5, capnp.Struct(ss).ToPtr())
	return ss, err
}

func (s Node) Directory() (Directory, error) {
	if capnp.Struct(s).Uint16(8) != 1 {
		panic("Which() != directory")
	}
	p, err := capnp.Struct(s).Ptr(5)
	return Directory(p.Struct()), err
}

func (s Node) HasDirectory() bool {
	if capnp.Struct(s).Uint16(8) != 1 {
		return false
	}
	return capnp.Struct(s).HasPtr(5)
}

func (s Node) SetDirectory(v Directory) error {
	capnp.Struct(s).SetUint16(8, 1)
	return capnp.Struct(s).SetPtr(5, capnp.Struct(v).ToPtr())
}

// NewDirectory sets the directory field to a newly
// allocated Directory struct, preferring placement in s's segment.
func (s Node) NewDirectory() (Directory, error) {
	capnp.Struct(s).SetUint16(8, 1)
	ss, err := NewDirectory(capnp.Struct(s).Segment())
	if err != nil {
		return Directory{}, err
	}
	err = capnp.Struct(s).SetPtr(5, capnp.Struct(ss).ToPtr())
	return ss, err
}

func (s Node) File() (File, error) {
	if capnp.Struct(s).Uint16(8) != 2 {
		panic("Which() != file")
	}
	p, err := capnp.Struct(s).Ptr(5)
	return File(p.Struct()), err
}

func (s Node) HasFile() bool {
	if capnp.Struct(s).Uint16(8) != 2 {
		return false
	}
	return capnp.Struct(s).HasPtr(5)
}

func (s Node) SetFile(v File) error {
	capnp.Struct(s).SetUint16(8, 2)
	return capnp.Struct(s).SetPtr(5, capnp.Struct(v).ToPtr())
}

// NewFile sets the file field to a newly
// allocated File struct, preferring placement in s's segment.
func (s Node) NewFile() (File, error) {
	capnp.Struct(s).SetUint16(8, 2)
	ss, err := NewFile(capnp.Struct(s).Segment())
	if err != nil {
		return File{}, err
	}
	err = capnp.Struct(s).SetPtr(5, capnp.Struct(ss).ToPtr())
	return ss, err
}

func (s Node) Ghost() (Ghost, error) {
	if capnp.Struct(s).Uint16(8) != 3 {
		panic("Which() != ghost")
	}
	p, err := capnp.Struct(s).Ptr(5)
	return Ghost(p.Struct()), err
}

func (s Node) HasGhost() bool {
	if capnp.Struct(s).Uint16(8) != 3 {
		return false
	}
	return capnp.Struct(s).HasPtr(5)
}

func (s Node) SetGhost(v Ghost) error {
	capnp.Struct(s).SetUint16(8, 3)
	return capnp.Struct(s).SetPtr(5, capnp.Struct(v).ToPtr())
}

// NewGhost sets the ghost field to a newly
// allocated Ghost struct, preferring placement in s's segment.
func (s Node) NewGhost() (Ghost, error) {
	capnp.Struct(s).SetUint16(8, 3)
	ss, err := NewGhost(capnp.Struct(s).Segment())
	if err != nil {
		return Ghost{}, err
	}
	err = capnp.Struct(s).SetPtr(5, capnp.Struct(ss).ToPtr())
	return ss, err
}

func (s Node) BackendHash() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(6)
	return []byte(p.Data()), err
}

func (s Node) HasBackendHash() bool {
	return capnp.Struct(s).HasPtr(6)
}

func (s Node) SetBackendHash(v []byte) error {
	return capnp.Struct(s).SetData(6, v)
}

// Node_List is a list of Node.
type Node_List = capnp.StructList[Node]

// NewNode creates a new list of Node.
func NewNode_List(s *capnp.Segment, sz int32) (Node_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 7}, sz)
	return capnp.StructList[Node](l), err
}

// Node_Future is a wrapper for a Node promised by a client call.
type Node_Future struct{ *capnp.Future }

func (f Node_Future) Struct() (Node, error) {
	p, err := f.Future.Ptr()
	return Node(p.Struct()), err
}
func (p Node_Future) Commit() Commit_Future {
	return Commit_Future{Future: p.Future.Field(5, nil)}
}
func (p Node_Future) Directory() Directory_Future {
	return Directory_Future{Future: p.Future.Field(5, nil)}
}
func (p Node_Future) File() File_Future {
	return File_Future{Future: p.Future.Field(5, nil)}
}
func (p Node_Future) Ghost() Ghost_Future {
	return Ghost_Future{Future: p.Future.Field(5, nil)}
}

const schema_9195d073cb5c5953 = "x\xda\xb4V]l\x14U\x18\xfd\xce\xbd\xb3;]R" +
	"\xec.\xb7$\x98\xd8\xec\x95`lI\xc5\x160B\x03" +
	"\x81j+\x05\x0b\xe9eK\x04\x03\xc6a\xf7\xb63a" +
	"w\xa6\x99\x19\xac5\x1a\xd0`\x02\x1a\x0cDL$)" +
	"\x11M\xfdK4\xf8LbLH4\x06_\x8c\x0f\x9a" +
	"\xf0\x88\x1a\xff\xa2\xcfB\x841w\x7fKmy\xf3q" +
	"\xce\xf7\xf5\xde\xf3\x9ds\xbf\xb3\xed\xbb\xc9\xb6[\xfd\xcb" +
	"\x1f\xb4\x88\xa9\xbeT:\xf9}\xc5\xec\xaf?t\x7f}" +
	"\x9c\xd4\x0a\xb0\xa4p\xe0\xe07\xd1\xb7o\x9d\xa5af" +
	"sX\xe2\x06\xae\x080[\x80\xe57ley\x10\x92" +
	"\x0b\xf9'\xa7\x9f\xfbk\xe5\xeb\x94[\x81V\x7f\x8a\xd9" +
	"D\xe2\x00\xbf&4\xb7\x85\xe6yq\x96O\x13\x92K" +
	"\x87\xc6\xfd\xaf\xc4\xc5\xd3\xe6\xf8\xf9\xedi\xd3\xfe'\xbf" +
	"*np[\xdc\xe0\xf9\x0d=\xd6S\xe6\xf4}\xfd\xa7" +
	"\x1e\xdd\xba\xf9\xc37\x16\xf6\xd7\x8eO]\x16N\xca\x16" +
	"N*/N\xa5.\x11\x92\x1foNL\x1d\xfb\xa3\xe7" +
	"\x83\x85\xecm\xdb\x82%\xeeO_\x16=i[\xf4\xa4" +
	"\xf3\x1b\x0e\xa5\x03FH\xe6~\x1a\xbd\xd61\xf7\xf7\x17" +
	"\xa4\xee\xc5<v+\xd36\x88\xc4\xcfm7\x09\xe2\xb7" +
	"6\xc3\x1c\xb3\xaf\x94\xfb\x0e\x8c^_\xc8\x84\x1b&\xc3" +
	"\x99\xebBel\xa12yq*\xf3\x0bmJ\x8a\xce" +
	"\x94?\xf5\xb0\x1f\xb0\x92\x8e\xd6U?\x06v\xb8A\x14" +
	"\xd3\x18\xa0,\xb0\xe4\x997\xdfQ\x9f\x7f\xff\xda\x97\xa4" +
	",\x86\xc1^\xa0\x9d\xa8\x1f\xdf!\xa9\xb6I\xcfO\x97" +
	"\xbc\xa2\x13\xebH\xc6\xae\x13KG\x16u\x18;\x9e/" +
	"\xfd\xa0\xa4\xe5\xb4\x13I'\x96\xb1\xebEr\xca\x89]" +
	"\x19\xf8Eh\"\xd5\xc9-\"\x0bD\xb9\x97\x9e&R" +
	"/r\xa8\x93\x0c@'\x0c\xf6\xea^\"u\x82C\x9d" +
	"a\xe8bI\x82N0\xa2\xdc\xe9\x01\"u\x92C\x9d" +
	"c\xe8\xe2\xb7\x0d\xcc\x89rgM\xf7\x19\x0e5\xcb\xd0" +
	"e\xdd2\xb0E\x94;\xbf\x96H\x9d\xe3P\x17\x19\x92" +
	"I\xc3v\xa7\x1f\x10/id\x88!Cup\xcc\x89" +
	"\x09.\xda\x89\xa1\x9d\xb0\xad\x18T*^\x8clKd" +
	"\x02\xb2\x84\xa4\xe4\x85\xba\x18\x07!a\x06\xd9\x96\xcc\xb5" +
	"j\xc7\x84W\xd6\xc8\xb6\xdeA\xfd\x8f\x16\x91w\xc8\x0b" +
	"\x87\xfd\x98\x873\x8b+|_U\xe1\x1c\xae&\x832" +
	"\xf2\xfc\xc9\xb2f\xb2q\xf5\x8c\xd4~\x1c\xce\x10T[" +
	"S\xbe\x1e3\xe5\x1a\x0e\xd5\xc7\x90k\xe8\xf7\x90\x01\xbb" +
	"9\xd4F\x86\x0e\xdf\xa9\xe8\xc6x\x1d\xae\x13\xb9XN" +
	"\x0c\xcb\x17g\xf7xu|Z\xc2~Y\xb7\x7f5\x92" +
	"Z\xa3\xf4x$\x1d\x19\xe9X\x06\x13\xb2\xe8:\xfe\xa4" +
	"y\x09\x81\xf4\x03\xbb\xa4#\"\xb5\xaa\xc9\xf4\xfc\xea\x96" +
	"\x1fM\xa6\x17\x8c\xa5os\xa89\x86\x1cc5\x9f\xdf" +
	"5\xe0,\x87\xfa\x88!\xc7y\xcd\xe5\xf7\xcdL\x179" +
	"\xd4'\x0c\xb0j\x16\x7f\xbc\x9eH\xcdq\xa8\xcf\x18\x90" +
	"\xc2\xbc=\xc9}\xba\x9e\x98]\x89&\x9b\xc6:Gc" +
	"7\x08\x9b\x9fSN\xa8\xfd\xb8!EG\x18\x04\xcd\x8f" +
	"\xbc\xe7\x97\xf4\xf3H\x11C\x8a\x90\xaf\xe8pR7\xb5" +
	"BC\xabmS\x03Oxe\xbd\xb8P\xab\xea.^" +
	"I\x06eY;\x13\xd2gf\x1d<_\xc6\xae\x96\xbb" +
	"\x87\x06w\x10\x91joj3l\x86\xdb\xce\xa1F[" +
	"K\xb0\xd3\xa80\xc4\xa1\xc6\x8c4\xf5\x15\xd8mD\x1c" +
	"\xe1P\xe3\x0c\x1d\x91\xf7B\xf317\x06\xaa\xcfg\x1f" +
	"\xd13\xff\xf1y>\xf7=Ai\x09\xeek\xea&\xef" +
	"B\xb2\xa7J:\x92\x96S[\xe7:\xff\x8a\x0e\x8f\x94" +
	"\xb5,9\x93\xc6\xf5\x89r\x10\x10Toc\x18\xf1\x00" +
	"\xd6\x12\x15$8\x0a\xbdhy-z\xb0\x8b\xa8\xd0m" +
	"\xf0\x8dh\xd9-\xfa\xf1\x18Q\xa1\xd7\xe0\x9b\xc0\x80\x9a" +
	"\xe1\xe2\x11\xac'*\xf4\x19x\x8bi\xb7x\xd5t\xb1" +
	"\x19\x87\x89\x0a\x9b\x0c>d\xf0\x94\xd5\x89\x14\x91\x18\xac" +
	"^\xbb\xc5\xe0#`\xe8J'I\xaa\x13i\x13{\x18" +
	" *l7\x95QS\xb1o\x9b\x8a\x09\xc4\x9d\xd8K" +
	"T\x181\x95qSi\xbbe*mDBUO\x1b" +
	"5\x95\xfd\xa6\x92\xf9\xc7T2Db_\x95\xd7\x98\xa9" +
	"\x1c4\xf7/Kwb\x99\x89\xf9*\xaf\xfd\x06/a" +
	"\xc1\xde%q\xa8\xf5\x88\x13\xb9D\xd4\xb0\xe5X%(" +
	"\x8d{\xad\x9e\xbcg4n\x86S1\xf0c\xed\xc7#" +
	"d\xcf[\xd9\x8e\xa3\x91\x0e\xff\x9f\xac\xcaW\xd3\x10\xd9" +
	"\xd6Ok\xfd\xb0\xc3N\xf1\x88\xf6Kw\x12Y:;" +
	"\xb6\xad\xab.\x8c\x89\xa8l\xcd\x99\x05\x19U3\xe5\xce" +
	"\x8c\x9a\xf6b\xb7\x95Q\xda)\xdd\xed\x9e\xa1\xeatv" +
	"\xb0T\x84v\xd7\x1f\xf0{H\x86\xeaB\xa4f\xa4\xd1" +
	"\xd3\xf1\xfcH\x06\xbe\x96A(+A\xa8\x9b\xc9\xea\xe9" +
	"\xc8`\x13\x9e]\xae\xa6V\xb6\xb9\x99\x8e\xa1y\x90C" +
	"\xb9\xad\xcd\xd4f3\x9f\xe5P\xe5y\x9b\xe9\xed\"R" +
	".\x87:aB\x8b\xd5B\xebe\x03\x1e\xaf\xfd4\xdd" +
	"m]\x93\xa2\xeb\x95K\xa1\xf6\xcd\xfb\xb8\x870\xc6\x81" +
	"l\xeb\x7f\x16\x82\x01\x1bO\"\xba[\xd3\xbf\x01\x00\x00" +
	"\xff\xff\xf1\x97$\x7f"

func init() {
	schemas.Register(schema_9195d073cb5c5953,
		0x80c828d7e89c12ea,
		0x8b15ee76774b1f9d,
		0x8da013c66e545daf,
		0x8ea7393d37893155,
		0xa629eb7f7066fae3,
		0xbff8a40fda4ce4a4,
		0xe24c59306c829c01)
}