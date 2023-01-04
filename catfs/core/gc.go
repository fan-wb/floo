package core

import (
	"floo/catfs/db"
	ie "floo/catfs/errors"
	n "floo/catfs/nodes"
)

// GarbageCollector implements a small mark & sweep garbage collector.
// It exists more for the sake of fault tolerance than it being an essential part of floo.
// This is different from the ipfs garbage collector.
type GarbageCollector struct {
	lkr      *Linker
	kv       db.Database
	notifier func(nd n.Node) bool
	markMap  map[string]struct{}
}

// NewGarbageCollector will return a new GC, operating on `lkr` and `kv`.
// It will call `kc` on every collected node.
func NewGarbageCollector(lkr *Linker, kv db.Database, kc func(nd n.Node) bool) *GarbageCollector {
	return &GarbageCollector{
		lkr:      lkr,
		kv:       kv,
		notifier: kc,
	}
}

func (gc *GarbageCollector) markMoveMap(key []string) error {
	keys, err := gc.kv.Keys(key...)
	if err != nil {
		return err
	}

	for _, key := range keys {
		data, err := gc.kv.Get(key...)
		if err != nil {
			return err
		}

		node, _, err := gc.lkr.parseMoveMappingLine(string(data))
		if err != nil {
			return err
		}

		if node != nil {
			gc.markMap[node.TreeHash().B58String()] = struct{}{}
		}
	}

	return nil
}

func (gc *GarbageCollector) mark(cmt *n.Commit, recursive bool) error {
	if cmt == nil {
		return nil
	}

	root, err := gc.lkr.DirectoryByHash(cmt.Root())
	if err != nil {
		return err
	}

	gc.markMap[cmt.TreeHash().B58String()] = struct{}{}
	err = n.Walk(gc.lkr, root, true, func(child n.Node) error {
		gc.markMap[child.TreeHash().B58String()] = struct{}{}
		return nil
	})

	if err != nil {
		return err
	}

	parent, err := cmt.Parent(gc.lkr)
	if err != nil {
		return err
	}

	if recursive && parent != nil {
		parentCmt, ok := parent.(*n.Commit)
		if !ok {
			return ie.ErrBadNode
		}

		return gc.mark(parentCmt, recursive)
	}

	return nil
}
