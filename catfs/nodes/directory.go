package nodes

import (
	h "floo/util/hashlib"
)

// Directory is a typical directory that may contain
// several other directories or files.
type Directory struct {
	Base

	size       uint64
	parentName string
	children   map[string]h.Hash
	contents   map[string]h.Hash
	order      []string
}
