package core

import (
	"floo/catfs/db"
	ie "floo/catfs/errors"
	n "floo/catfs/nodes"
	h "floo/util/hashlib"
)

// GarbageCollector implements a small mark & sweep garbage collector.
// It exists more for the sake of fault tolerance than being an essential part of floo.
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

func (gc *GarbageCollector) sweep(prefix []string) (int, error) {
	removed := 0

	return removed, gc.lkr.AtomicWithBatch(func(batch db.Batch) (bool, error) {
		keys, err := gc.kv.Keys(prefix...)
		if err != nil {
			return hintRollback(err)
		}

		for _, key := range keys {
			b58Hash := key[len(key)-1]
			if _, ok := gc.markMap[b58Hash]; ok {
				continue
			}

			hash, err := h.FromB58String(b58Hash)
			if err != nil {
				return hintRollback(err)
			}

			node, err := gc.lkr.NodeByHash(hash)
			if err != nil {
				return hintRollback(err)
			}

			if node == nil {
				continue
			}

			// Allow the gc caller to check if he really
			// wants to delete this node.
			if gc.notifier != nil && !gc.notifier(node) {
				continue
			}

			// Actually get rid of the node:
			gc.lkr.MemIndexPurge(node)

			batch.Erase(key...)
			removed++
		}

		return false, nil
	})
}

func (gc *GarbageCollector) findAllMoveLocations(head *n.Commit) ([][]string, error) {
	locations := [][]string{
		{"stage", "moves"},
	}

	for {
		parent, err := head.Parent(gc.lkr)
		if err != nil {
			return nil, err
		}

		if parent == nil {
			break
		}

		parentCmt, ok := parent.(*n.Commit)
		if !ok {
			return nil, ie.ErrBadNode
		}

		head = parentCmt
		location := []string{"moves", head.TreeHash().B58String()}
		locations = append(locations, location)
	}

	return locations, nil
}
