package core

import (
	"capnproto.org/go/capnp/v3"
	"encoding/binary"
	"floo/catfs/db"
	ie "floo/catfs/errors"
	n "floo/catfs/nodes"
	h "floo/util/hashlib"
	"floo/util/trie"
	"fmt"
	e "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// Linker implements the basic logic of floo's data model
// It uses an underlying key/value database to
// store a Merkle-DAG with versioned metadata,
// similar to what git does internally.
type Linker struct {
	kv db.Database

	// root of the filesystem
	root *n.Directory

	// Path lookup trie
	ptrie *trie.Node

	// B58Hash to node
	index map[string]n.Node

	// UID to node
	inodeIndex map[uint64]n.Node

	// Cache for the linker owner.
	owner string
}

// NewLinker returns a new lkr, ready to use. It assumes the key value store
// is working and does no check on this.
func NewLinker(kv db.Database) *Linker {
	lkr := &Linker{kv: kv}
	lkr.MemIndexClear()
	return lkr
}

// MemIndexAdd adds `nd` to the in memory index.
func (lkr *Linker) MemIndexAdd(nd n.Node, updatePathIndex bool) {
	lkr.index[nd.TreeHash().B58String()] = nd
	lkr.inodeIndex[nd.Inode()] = nd

	if updatePathIndex {
		path := nd.Path()
		if nd.Type() == n.NodeTypeDirectory {
			path = appendDot(path)
		}
		lkr.ptrie.InsertWithData(path, nd)
	}
}

// MemIndexSwap updates an entry of the in memory index, by deleting
// the old entry referenced by oldHash (maybe nil). This is necessary
// to ensure that old hashes do not resolve to the new, updated instance.
// If the old instance is needed, it will be loaded as new instance.
// You should not need to call this function, except when implementing own Nodes.
func (lkr *Linker) MemIndexSwap(nd n.Node, oldHash h.Hash, updatePathIndex bool) {
	if oldHash != nil {
		delete(lkr.index, oldHash.B58String())
	}

	lkr.MemIndexAdd(nd, updatePathIndex)
}

// MemSetRoot sets the current root, but does not store it yet. It's supposed
// to be called after in-memory modifications. Only implementors of new Nodes
// might need to call this function.
func (lkr *Linker) MemSetRoot(root *n.Directory) {
	if lkr.root != nil {
		lkr.MemIndexSwap(root, lkr.root.TreeHash(), true)
	} else {
		lkr.MemIndexAdd(root, true)
	}

	lkr.root = root
}

// MemIndexPurge removes `nd` from the memory index.
func (lkr *Linker) MemIndexPurge(nd n.Node) {
	delete(lkr.inodeIndex, nd.Inode())
	delete(lkr.index, nd.TreeHash().B58String())
	lkr.ptrie.Lookup(nd.Path()).Remove()
}

// MemIndexClear resets the memory index to zero.
// This should not be called mid-flight in operations,
// but should be okay to call between atomic operations.
func (lkr *Linker) MemIndexClear() {
	lkr.ptrie = trie.NewNode()
	lkr.index = make(map[string]n.Node)
	lkr.inodeIndex = make(map[uint64]n.Node)
	lkr.root = nil
}

//////////////////////////
// COMMON NODE HANDLING //
//////////////////////////

// NextInode returns a unique identifier, used to identify a single node. You
// should not need to call this function, except when implementing own nodes.
func (lkr *Linker) NextInode() uint64 {
	nodeCount, err := lkr.kv.Get("stats", "max-inode")
	if err != nil && err != db.ErrNoSuchKey {
		return 0
	}

	// nodeCount might be nil on startup:
	cnt := uint64(1)
	if nodeCount != nil {
		cnt = binary.BigEndian.Uint64(nodeCount) + 1
	}

	cntBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(cntBuf, cnt)

	err = lkr.AtomicWithBatch(func(batch db.Batch) (bool, error) {
		batch.Put(cntBuf, "stats", "max-inode")
		return false, nil
	})

	if err != nil {
		return 0
	}

	return cnt
}

// FilesByContents checks what files are associated with the content hashes in
// `contents`. It returns a map of content hash b58 to file. This method is
// quite heavy and should not be used in loops. There is room for optimizations.
func (lkr *Linker) FilesByContents(contents []h.Hash) (map[string]*n.File, error) {
	keys, err := lkr.kv.Keys()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*n.File)
	for _, key := range keys {
		// Filter non-node storage:
		fullKey := strings.Join(key, "/")
		if !strings.HasPrefix(fullKey, "objects") &&
			!strings.HasPrefix(fullKey, "stage/objects") {
			continue
		}

		data, err := lkr.kv.Get(key...)
		if err != nil {
			return nil, err
		}

		nd, err := n.UnmarshalNode(data)
		if err != nil {
			return nil, err
		}

		if nd.Type() != n.NodeTypeFile {
			continue
		}

		file, ok := nd.(*n.File)
		if !ok {
			return nil, ie.ErrBadNode
		}

		for _, content := range contents {
			if content.Equal(file.BackendHash()) {
				result[content.B58String()] = file
			}
		}
	}

	return result, nil
}

// loadNode loads an individual object by its hash from the object store.
// It will return nil if the hash is not there.
func (lkr *Linker) loadNode(hash h.Hash) (n.Node, error) {
	var data []byte
	var err error

	b58hash := hash.B58String()

	// First look in the stage:
	loadableBuckets := [][]string{
		{"stage", "objects", b58hash},
		{"objects", b58hash},
	}

	for _, bucketPath := range loadableBuckets {
		data, err = lkr.kv.Get(bucketPath...)
		if err != nil && err != db.ErrNoSuchKey {
			return nil, err
		}

		if data != nil {
			return n.UnmarshalNode(data)
		}
	}

	// Damn, no hash found:
	return nil, nil
}

// NodeByHash returns the node identified by hash.
// If no such hash could be found, nil is returned.
func (lkr *Linker) NodeByHash(hash h.Hash) (n.Node, error) {
	// Check if we have this node in the memory cache already:
	b58Hash := hash.B58String()
	if cachedNode, ok := lkr.index[b58Hash]; ok {
		return cachedNode, nil
	}

	// Node was not in the cache, load directly from kv.
	nd, err := lkr.loadNode(hash)
	if err != nil {
		return nil, err
	}

	if nd == nil {
		return nil, nil
	}

	lkr.MemIndexAdd(nd, false)
	return nd, nil
}

// NodeByInode resolves a node by its unique ID.
// It will return nil if no corresponding node was found.
func (lkr *Linker) NodeByInode(uid uint64) (n.Node, error) {
	b58Hash, err := lkr.kv.Get("inode", strconv.FormatUint(uid, 10))
	if err != nil && err != db.ErrNoSuchKey {
		return nil, err
	}

	hash, err := h.FromB58String(string(b58Hash))
	if err != nil {
		return nil, err
	}

	return lkr.NodeByHash(hash)
}

func appendDot(path string) string {
	// path.Join() calls path.Clean() which in turn
	// removes the '.' at the end when trying to join that.
	// But since we use the dot to mark directories we shouldn't do that.
	if strings.HasSuffix(path, "/") {
		return path + "."
	}

	return path + "/."
}

// ResolveNode resolves a path to a hash and resolves the corresponding node by
// calling NodeByHash(). If no node could be resolved, nil is returned.
// It does not matter if the node was deleted in the meantime. If so,
// a Ghost node is returned which stores the last known state.
func (lkr *Linker) ResolveNode(nodePath string) (n.Node, error) {
	// Check if it's cached already:
	trieNode := lkr.ptrie.Lookup(nodePath)
	if trieNode != nil && trieNode.Data != nil {
		return trieNode.Data.(n.Node), nil
	}

	fullPaths := [][]string{
		{"stage", "tree", nodePath},
		{"tree", nodePath},
	}

	for _, fullPath := range fullPaths {
		b58Hash, err := lkr.kv.Get(fullPath...)
		if err != nil && err != db.ErrNoSuchKey {
			return nil, e.Wrapf(err, "db-lookup")
		}

		if err == db.ErrNoSuchKey {
			continue
		}

		bhash, err := h.FromB58String(string(b58Hash))
		if err != nil {
			return nil, err
		}

		if bhash != nil {
			return lkr.NodeByHash(h.Hash(bhash))
		}
	}

	// Return nil if nothing found:
	return nil, nil
}

// StageNode inserts a modified node to the staging area, making sure the
// modification is persistent and part of the staging commit. All parent
// directories of the node in question will be staged automatically. If there
// was no modification it will be a (quite expensive) NOOP.
func (lkr *Linker) StageNode(nd n.Node) error {
	return lkr.AtomicWithBatch(func(batch db.Batch) (bool, error) {
		if err := lkr.stageNodeRecursive(batch, nd); err != nil {
			return true, e.Wrapf(err, "recursive stage")
		}

		// Update the staging commit's root hash:
		status, err := lkr.Status()
		if err != nil {
			return true, fmt.Errorf("failed to retrieve status: %v", err)
		}

		root, err := lkr.Root()
		if err != nil {
			return true, err
		}

		status.SetModTime(time.Now())
		status.SetRoot(root.TreeHash())
		lkr.MemSetRoot(root)
		return hintRollback(lkr.saveStatus(status))
	})
}

func (lkr *Linker) stageNodeRecursive(batch db.Batch, nd n.Node) error {
	if nd.Type() == n.NodeTypeCommit {
		return fmt.Errorf("bug: commits cannot be staged; use MakeCommit()")
	}

	data, err := n.MarshalNode(nd)
	if err != nil {
		return e.Wrapf(err, "marshal")
	}

	b58Hash := nd.TreeHash().B58String()
	batch.Put(data, "stage", "objects", b58Hash)

	uidKey := strconv.FormatUint(nd.Inode(), 10)
	batch.Put([]byte(nd.TreeHash().B58String()), "inode", uidKey)

	hashPath := []string{"stage", "tree", nd.Path()}
	if nd.Type() == n.NodeTypeDirectory {
		hashPath = append(hashPath, ".")
	}

	batch.Put([]byte(b58Hash), hashPath...)

	// Remember/Update this node in the cache if it's not yet there:
	lkr.MemIndexAdd(nd, true)

	// We need to save parent directories too, in case the hash changed:
	// Note that this will create many pointless directories in staging.
	// That's okay since we garbage-collect it every few seconds
	// on a higher layer.
	if nd.Path() == "/" {
		// Can't go any higher. Save this dir as new virtual root.
		root, ok := nd.(*n.Directory)
		if !ok {
			return ie.ErrBadNode
		}

		lkr.MemSetRoot(root)
		return nil
	}

	par, err := lkr.ResolveDirectory(path.Dir(nd.Path()))
	if err != nil {
		return e.Wrapf(err, "resolve")
	}

	if par != nil {
		if err := lkr.stageNodeRecursive(batch, par); err != nil {
			return err
		}
	}

	return nil
}

// CommitByIndex returns the commit referenced by `index`.
// `0` will return the very first commit. Negative numbers will yield
// a ErrNoSuchKey error.
func (lkr *Linker) CommitByIndex(index int64) (*n.Commit, error) {
	status, err := lkr.Status()
	if err != nil {
		return nil, err
	}

	if index < 0 {
		// Interpret an index of -n as curr-(n+1),
		// so that -1 means "curr".
		index = status.Index() + index + 1
	}

	b58Hash, err := lkr.kv.Get("index", strconv.FormatInt(index, 10))
	if err != nil && err != db.ErrNoSuchKey {
		return nil, err
	}

	// Special case: status is not in the index bucket.
	// Do a separate check for it.
	if err == db.ErrNoSuchKey {
		if status.Index() == index {
			return status, nil
		}

		return nil, nil
	}

	hash, err := h.FromB58String(string(b58Hash))
	if err != nil {
		return nil, err
	}

	return lkr.CommitByHash(hash)
}

/////////////////////
// COMMIT HANDLING //
/////////////////////

////////////////////////
// REFERENCE HANDLING //
////////////////////////

// ResolveRef resolves the hash associated with `refName`. If the ref could not
// be resolved, ErrNoSuchRef is returned. Typically, Node will be a Commit.
// But there are no technical restrictions on which node typ to use.
// NOTE: ResolveRef("HEAD") != ResolveRef("head") due to case.
func (lkr *Linker) ResolveRef(refName string) (n.Node, error) {
	origRefName := refName

	nUps := 0
	for idx := len(refName) - 1; idx >= 0; idx-- {
		if refName[idx] == '^' {
			nUps++
		} else {
			break
		}
	}

	// Strip the ^s:
	refName = refName[:len(refName)-nUps]

	// Special case: the status commit is not part of the normal object store.
	// Still make it able to resolve it by its refName "curr".
	if refName == "curr" || refName == "status" {
		return lkr.Status()
	}

	b58Hash, err := lkr.kv.Get("refs", refName)
	if err != nil && err != db.ErrNoSuchKey {
		return nil, err
	}

	if len(b58Hash) == 0 {
		// Try to interpret the refName as b58hash directly.
		// This path will hit when passing a commit hash directly
		// as `refName` to this method.
		b58Hash = []byte(refName)
	}

	hash, err := h.FromB58String(string(b58Hash))
	if err != nil {
		// Could not parse hash, so it's probably none.
		return nil, ie.ErrNoSuchRef(refName)
	}

	status, err := lkr.Status()
	if err != nil {
		return nil, err
	}

	// Special case: Allow the resolving of `curr`
	// by using its status hash and check it explicitly.
	var nd n.Node
	if status.TreeHash().Equal(hash) {
		nd = status
	} else {
		nd, err = lkr.NodeByHash(h.Hash(hash))
		if err != nil {
			return nil, err
		}
	}

	if nd == nil {
		return nil, ie.ErrNoSuchRef(refName)
	}

	// Possibly advance a few commits until we hit the one
	// the user required.
	cmt, ok := nd.(*n.Commit)
	if ok {
		for i := 0; i < nUps; i++ {
			parentNd, err := cmt.Parent(lkr)
			if err != nil {
				return nil, err
			}

			if parentNd == nil {
				log.Warningf("ref `%s` is too far back; stopping at `init`", origRefName)
				break
			}

			parentCmt, ok := parentNd.(*n.Commit)
			if !ok {
				break
			}

			cmt = parentCmt
		}

		nd = cmt
	}

	return nd, nil
}

// SaveRef stores a reference to `nd` persistently. The caller is responsible
// to ensure that the node is already in the blockstore, otherwise it won't be
// resolvable.
func (lkr *Linker) SaveRef(refName string, nd n.Node) error {
	refName = strings.ToLower(refName)
	return lkr.AtomicWithBatch(func(batch db.Batch) (bool, error) {
		batch.Put([]byte(nd.TreeHash().B58String()), "refs", refName)
		return false, nil
	})
}

// ListRefs lists all currently known refs.
func (lkr *Linker) ListRefs() ([]string, error) {
	refs := []string{}
	keys, err := lkr.kv.Keys("refs")
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		if len(key) <= 1 {
			continue
		}

		refs = append(refs, key[1])
	}

	return refs, nil
}

// RemoveRef removes the ref named `refName`.
func (lkr *Linker) RemoveRef(refName string) error {
	return lkr.AtomicWithBatch(func(batch db.Batch) (bool, error) {
		batch.Erase("refs", refName)
		return false, nil
	})
}

// Head is just a shortcut for ResolveRef("HEAD").
func (lkr *Linker) Head() (*n.Commit, error) {
	nd, err := lkr.ResolveRef("head")
	if err != nil {
		return nil, err
	}

	cmt, ok := nd.(*n.Commit)
	if !ok {
		return nil, fmt.Errorf("uh-oh, HEAD is not a Commit... %v", nd)
	}

	return cmt, nil
}

// Root returns the current root directory of CURR.
// It is never nil when err is nil.
func (lkr *Linker) Root() (*n.Directory, error) {
	if lkr.root != nil {
		return lkr.root, nil
	}

	status, err := lkr.Status()
	if err != nil {
		return nil, err
	}

	rootNd, err := lkr.DirectoryByHash(status.Root())
	if err != nil {
		return nil, err
	}

	lkr.MemSetRoot(rootNd)
	return rootNd, nil
}

// Status returns the current staging commit.
// It is never nil, unless err is nil.
func (lkr *Linker) Status() (*n.Commit, error) {
	var cmt *n.Commit
	var err error

	return cmt, lkr.AtomicWithBatch(func(batch db.Batch) (bool, error) {
		cmt, err = lkr.status(batch)
		return hintRollback(err)
	})
}

func (lkr *Linker) status(batch db.Batch) (cmt *n.Commit, err error) {
	cmt, err = lkr.loadStatus()
	if err != nil {
		return nil, err
	}

	if cmt != nil {
		return cmt, nil
	}

	// Shoot, no commit exists yet.
	// We need to create an initial one.
	cmt, err = n.NewEmptyCommit(lkr.NextInode(), 0)
	if err != nil {
		return nil, err
	}

	// Set up a new commit and set root from last HEAD or new one.
	head, err := lkr.Head()
	if err != nil && !ie.IsErrNoSuchRef(err) {
		return nil, err
	}

	var rootHash h.Hash

	if ie.IsErrNoSuchRef(err) {
		// There probably wasn't a HEAD yet.
		if root, err := lkr.ResolveDirectory("/"); err == nil && root != nil {
			rootHash = root.TreeHash()
		} else {
			// No root directory then. Create a shiny new one and stage it.
			inode := lkr.NextInode()
			newRoot, err := n.NewEmptyDirectory(lkr, nil, "/", lkr.owner, inode)
			if err != nil {
				return nil, err
			}

			// Can't call StageNode(), since that would call Status(),
			// causing and endless loop of grief and doom.
			if err := lkr.stageNodeRecursive(batch, newRoot); err != nil {
				return nil, err
			}

			rootHash = newRoot.TreeHash()
		}
	} else {
		if err := cmt.SetParent(lkr, head); err != nil {
			return nil, err
		}

		rootHash = head.Root()
	}

	cmt.SetRoot(rootHash)

	if err := lkr.saveStatus(cmt); err != nil {
		return nil, err
	}

	return cmt, nil
}

func (lkr *Linker) loadStatus() (*n.Commit, error) {
	data, err := lkr.kv.Get("stage", "STATUS")
	if err != nil && err != db.ErrNoSuchKey {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	msg, err := capnp.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	// It's there already. Just unmarshal it.
	cmt := &n.Commit{}
	if err := cmt.FromCapnp(msg); err != nil {
		return nil, err
	}

	return cmt, nil
}

// saveStatus copies cmt to stage/STATUS.
func (lkr *Linker) saveStatus(cmt *n.Commit) error {
	return lkr.AtomicWithBatch(func(batch db.Batch) (bool, error) {
		head, err := lkr.Head()
		if err != nil && !ie.IsErrNoSuchRef(err) {
			return hintRollback(err)
		}

		if head != nil {
			if err := cmt.SetParent(lkr, head); err != nil {
				return hintRollback(err)
			}
		}

		if err := cmt.BoxCommit(n.AuthorOfStage, ""); err != nil {
			return hintRollback(err)
		}

		data, err := n.MarshalNode(cmt)
		if err != nil {
			return hintRollback(err)
		}

		inode := strconv.FormatUint(cmt.Inode(), 10)
		batch.Put(data, "stage", "STATUS")
		batch.Put([]byte(cmt.TreeHash().B58String()), "inode", inode)
		return hintRollback(lkr.SaveRef("CURR", cmt))
	})
}

/////////////////////////////////
// CONVENIENT ACCESS FUNCTIONS //
/////////////////////////////////

// LookupNode takes the root node and tries to resolve the path from there.
// Deleted paths are recognized in contrast to ResolveNode.
// If a path does not exist NoSuchFile is returned.
func (lkr *Linker) LookupNode(repoPath string) (n.Node, error) {
	root, err := lkr.Root()
	if err != nil {
		return nil, err
	}

	return root.Lookup(lkr, repoPath)
}

// DirectoryByHash calls NodeByHash and attempts to convert
// it to a Directory as convenience.
func (lkr *Linker) DirectoryByHash(hash h.Hash) (*n.Directory, error) {
	nd, err := lkr.NodeByHash(hash)
	if err != nil {
		return nil, err
	}

	if nd == nil {
		return nil, nil
	}

	dir, ok := nd.(*n.Directory)
	if !ok {
		return nil, ie.ErrBadNode
	}

	return dir, nil
}

// CommitByHash lookups a commit by its hash.
// If the commit could not be found, nil is returned.
func (lkr *Linker) CommitByHash(hash h.Hash) (*n.Commit, error) {
	nd, err := lkr.NodeByHash(hash)
	if err != nil {
		return nil, err
	}

	if nd == nil {
		return nil, nil
	}

	cmt, ok := nd.(*n.Commit)
	if !ok {
		return nil, ie.ErrBadNode
	}

	return cmt, nil
}

// helper to return errors that should trigger a rollback in AtomicWithBatch()
func hintRollback(err error) (bool, error) {
	if err != nil {
		return true, err
	}

	return false, nil
}

// Atomic is like AtomicWithBatch but does not require using a batch.
// Use this for read-only operations. It's only syntactic sugar though.
func (lkr *Linker) Atomic(fn func() (bool, error)) (err error) {
	return lkr.AtomicWithBatch(func(batch db.Batch) (bool, error) {
		return fn()
	})
}

// AtomicWithBatch will execute `fn` in one transaction.
// If anything goes wrong (i.e. `fn` returns an error)
func (lkr *Linker) AtomicWithBatch(fn func(batch db.Batch) (bool, error)) (err error) {
	batch := lkr.kv.Batch()

	// A panicking program should not leave the persistent linker state
	// inconsistent. This is really a last defence against all odds.
	defer func() {
		if r := recover(); r != nil {
			batch.Rollback()
			lkr.MemIndexClear()
			err = fmt.Errorf("panic rollback: %v; stack: %s", r, string(debug.Stack()))
		}
	}()

	needRollback, err := fn(batch)
	if needRollback && err != nil {
		hadWrites := batch.HaveWrites()
		batch.Rollback()

		// Only clear the whole index if something was written.
		// Also, this prevents the slightly misleading log message below
		// in case of read-only operations.
		if hadWrites {
			// clearing the mem index will cause it to be read freshly from disk with the old state.
			lkr.MemIndexClear()
			log.Warningf("rolled back due to error: %v %s", err, debug.Stack())
		}

		return err
	}

	// Attempt to write it to disk.
	// If that fails we're better off deleting our internal cache.
	// so memory and disk is in sync.
	if flushErr := batch.Flush(); flushErr != nil {
		lkr.MemIndexClear()
		log.Warningf("flush to db failed, resetting mem index: %v", flushErr)
	}

	return err
}

// ResolveDirectory calls ResolveNode and converts the result to a Directory.
// This only accesses nodes from the filesystem and does not differentiate
// between ghosts and living nodes.
func (lkr *Linker) ResolveDirectory(dirPath string) (*n.Directory, error) {
	nd, err := lkr.ResolveNode(appendDot(path.Clean(dirPath)))
	if err != nil {
		return nil, err
	}

	if nd == nil {
		return nil, nil
	}

	dir, ok := nd.(*n.Directory)
	if !ok {
		return nil, ie.ErrBadNode
	}

	return dir, nil
}
