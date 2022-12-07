package nodes

import (
	h "floo/util/hashlib"
)

const (
	// AuthorOfStage is the Person that is displayed for the stage commit.
	// Currently, this is just an empty hash Person that will be set later.
	AuthorOfStage = "unknown"
)

// Commit groups a set of changes
type Commit struct {
	Base

	// commit message (might be auto-generated)
	message string

	// author is the id of the committer.
	author string

	// root is the tree hash of the root directory
	root h.Hash

	// parent hash (only nil for initial commit)
	parent h.Hash

	// index of the commit (first is 0, second 1 and so on)
	index int64

	merge struct {
		// with indicates with which person we merged.
		with string

		// head is a reference to the commit we merged with on
		// the remote side.
		head h.Hash
	}
}
