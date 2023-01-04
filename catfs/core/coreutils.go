package core

import (
	ie "floo/catfs/errors"
	n "floo/catfs/nodes"
	h "floo/util/hashlib"
	"fmt"
	e "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"path"
	"strings"
	"time"
)

// mkdirParents takes the dirname of repoPath and makes sure all intermediate
// directories are created. The last directory will be returned.

// If any directory exist already, it will not be touched.
// You can also think of it as mkdir -p.
func mkdirParents(lkr *Linker, repoPath string) (*n.Directory, error) {
	repoPath = path.Clean(repoPath)

	elems := strings.Split(repoPath, "/")
	for idx := 0; idx < len(elems)-1; idx++ {
		dirname := strings.Join(elems[:idx+1], "/")
		if dirname == "" {
			dirname = "/"
		}

		dir, err := Mkdir(lkr, dirname, false)
		if err != nil {
			return nil, err
		}

		// Return it, if it's the last path component:
		if idx+1 == len(elems)-1 {
			return dir, nil
		}
	}

	return nil, fmt.Errorf("empty path given")
}

// Mkdir creates the directory at repoPath and any intermediate directories if
// createParents is true. It will fail if there is already a file at `repoPath`
// and it is not a directory.
func Mkdir(lkr *Linker, repoPath string, createParents bool) (dir *n.Directory, err error) {
	dirname, basename := path.Split(repoPath)

	// Take special care of the root node:
	if basename == "" {
		return lkr.Root()
	}

	// Check if the parent exists:
	parent, lerr := lkr.LookupDirectory(dirname)
	if lerr != nil && !ie.IsNoSuchFileError(lerr) {
		err = e.Wrap(lerr, "dirname lookup failed")
		return
	}

	err = lkr.Atomic(func() (bool, error) {
		// If it's nil, we might need to create it:
		if parent == nil {
			if !createParents {
				return false, ie.NoSuchFile(dirname)
			}

			parent, err = mkdirParents(lkr, repoPath)
			if err != nil {
				return true, err
			}
		}

		child, err := parent.Child(lkr, basename)
		if err != nil {
			return true, err
		}

		if child != nil {
			switch child.Type() {
			case n.NodeTypeDirectory:
				// Nothing to do really. Return the old child.
				dir = child.(*n.Directory)
				return false, nil
			case n.NodeTypeFile:
				return true, fmt.Errorf("`%s` exists and is a file", repoPath)
			case n.NodeTypeGhost:
				// Remove the ghost and continue with adding:
				if err := parent.RemoveChild(lkr, child); err != nil {
					return true, err
				}
			default:
				return true, ie.ErrBadNode
			}
		}

		// Create it then!
		dir, err = n.NewEmptyDirectory(lkr, parent, basename, lkr.owner, lkr.NextInode())
		if err != nil {
			return true, err
		}

		if err := lkr.StageNode(dir); err != nil {
			return true, e.Wrapf(err, "stage dir")
		}

		log.Debugf("mkdir: %s", dirname)
		return false, nil
	})

	return
}

// Stage adds a file to floo's DAG.
func Stage(lkr *Linker, repoPath string, contentHash, backendHash h.Hash, size uint64, key []byte) (file *n.File, err error) {
	node, lerr := lkr.LookupNode(repoPath)
	if lerr != nil && !ie.IsNoSuchFileError(lerr) {
		err = lerr
		return
	}

	err = lkr.Atomic(func() (bool, error) {
		if node != nil {
			if node.Type() == n.NodeTypeGhost {
				ghostParent, err := n.ParentDirectory(lkr, node)
				if err != nil {
					return true, err
				}

				if ghostParent == nil {
					return true, fmt.Errorf(
						"bug: %s has no parent. Is root a ghost?",
						node.Path(),
					)
				}

				if err := ghostParent.RemoveChild(lkr, node); err != nil {
					return true, err
				}

				// Act like there was no previous node.
				// New node will have a different Inode.
				file = nil
			} else {
				var ok bool
				file, ok = node.(*n.File)
				if !ok {
					return true, ie.ErrBadNode
				}
			}
		}

		needRemove := false
		if file != nil {
			// We know this file already.
			log.WithFields(log.Fields{"file": repoPath}).Info("File exists; modifying.")
			needRemove = true

			if file.BackendHash().Equal(backendHash) {
				log.Debugf("Hash was not modified. Not doing any update.")
				return false, nil
			}
		} else {
			parent, err := mkdirParents(lkr, repoPath)
			if err != nil {
				return true, err
			}

			// Create a new file at specified path:
			file = n.NewEmptyFile(parent, path.Base(repoPath), lkr.owner, lkr.NextInode())
		}

		parentDir, err := n.ParentDirectory(lkr, file)
		if err != nil {
			return true, err
		}

		if parentDir == nil {
			return true, fmt.Errorf("%s has no parent yet (BUG)", repoPath)
		}

		if needRemove {
			// Remove the child before changing the hash:
			if err := parentDir.RemoveChild(lkr, file); err != nil {
				return true, err
			}
		}

		file.SetSize(size)
		file.SetModTime(time.Now())
		file.SetContent(lkr, contentHash)
		file.SetBackend(lkr, backendHash)
		file.SetKey(key)
		file.SetUser(lkr.owner)

		// Add it again when the hash was changed.
		log.Debugf("adding %s (%v)", file.Path(), file.BackendHash())
		if err := parentDir.Add(lkr, file); err != nil {
			return true, err
		}

		if err := lkr.StageNode(file); err != nil {
			return true, err
		}

		return false, nil
	})

	return
}
