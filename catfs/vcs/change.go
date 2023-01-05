package vcs

import (
	n "floo/catfs/nodes"
	"fmt"
	"strings"
)

const (
	// ChangeTypeNone means that a node did not change (compared to HEAD)
	ChangeTypeNone = ChangeType(0)
	// ChangeTypeAdd says that the node was initially added after HEAD.
	ChangeTypeAdd = ChangeType(1 << iota)
	// ChangeTypeModify says that the node was modified after HEAD
	ChangeTypeModify
	// ChangeTypeMove says that the node was moved after HEAD.
	// Note that Move and Modify may happen at the same time.
	ChangeTypeMove
	// ChangeTypeRemove says that the node was removed after HEAD.
	ChangeTypeRemove
)

// ChangeType is a mask of possible state change events.
type ChangeType uint8

// String will convert a ChangeType to a human-readable form
func (ct ChangeType) String() string {
	v := []string{}

	if ct&ChangeTypeAdd != 0 {
		v = append(v, "added")
	}
	if ct&ChangeTypeModify != 0 {
		v = append(v, "modified")
	}
	if ct&ChangeTypeMove != 0 {
		v = append(v, "moved")
	}
	if ct&ChangeTypeRemove != 0 {
		v = append(v, "removed")
	}

	if len(v) == 0 {
		return "none"
	}

	return strings.Join(v, "|")
}

// IsCompatible checks if two change masks are compatible.
// Changes are compatible when they can be both applied
// without losing any content. We may lose metadata though,
// e.g. when one side was moved, but the other removed:
// Here the remove would win and no move is counted.
func (ct ChangeType) IsCompatible(ot ChangeType) bool {
	modifyMask := ChangeTypeAdd | ChangeTypeModify
	return ct&modifyMask == 0 || ot&modifyMask == 0
}

// Change represents a single change of a node between two commits.
type Change struct {
	// Mask is a bitmask of changes that were made.
	// It describes the change that was made between `Next` to `Head`
	// and which is part of `Head`.
	Mask ChangeType

	// Head is the commit that was the current HEAD when this change happened.
	// Note that this is NOT the commit that contains the change, but the commit before.
	Head *n.Commit

	// Next is the commit that comes before `Head`.
	Next *n.Commit

	// Curr is the node with the attributes at a specific state
	Curr n.ModNode

	// MovedTo is only filled for ghosts that were the source
	// of a move. It's the path of the node it was moved to.
	MovedTo string

	// WasPreviouslyAt points to the place `Curr` was at
	// before a move. On changes without a move this is empty.
	WasPreviouslyAt string
}

func (ch *Change) String() string {
	movedTo := ""
	if len(ch.MovedTo) != 0 {
		movedTo = fmt.Sprintf(" (now %s)", ch.MovedTo)
	}

	prevAt := ""
	if len(ch.WasPreviouslyAt) != 0 {
		prevAt = fmt.Sprintf(" (was %s)", ch.WasPreviouslyAt)
	}

	return fmt.Sprintf("<%s:%s%s%s>", ch.Curr.Path(), ch.Mask, prevAt, movedTo)
}
