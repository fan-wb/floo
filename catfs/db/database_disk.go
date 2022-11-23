package db

import (
	"floo/util"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// DiskDatabase is a database that uses a filesystem as storage
// Each bucket is a directory, Leaf keys are single files
// The exported form of the database is simply a gzipped .tar of the directory
type DiskDatabase struct {
	basePath string
	refs     int64
	cache    map[string][]byte
	ops      []func() error
	deletes  map[string]struct{}
}

// NewDiskDatabase creates a new database at `basePath`.
func NewDiskDatabase(basePath string) (*DiskDatabase, error) {
	return &DiskDatabase{
		basePath: basePath,
		cache:    make(map[string][]byte),
		deletes:  make(map[string]struct{}),
	}, nil
}

// Put stores a new `val` under `key` at bucket.
func (db *DiskDatabase) Put(val []byte, key ...string) {

	db.ops = append(db.ops, func() error {
		filePath := filepath.Join(db.basePath, fixDirectoryKeys(key))

		// remove any non-directory parent to enable the nesting
		parentDir := filepath.Dir(filePath)
		if err := removeNonDirs(parentDir); err != nil {
			return err
		}

		if err := os.Mkdir(parentDir, 0700); err != nil {
			return err
		}

		// it is allowed to set key over an existing one
		// i.e. set "a/b" over "a/b/c".
		//This requires us to potentially delete nested directories (c).
		info, err := os.Stat(parentDir)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if info != nil && info.IsDir() {
			if err := os.RemoveAll(filePath); err != nil {
				return err
			}
		}
		return os.WriteFile(filePath, val, 0600)
	})
}

// Clear removes all keys below and including `key`
func (db *DiskDatabase) Clear(key ...string) error {

	// cache the real modification for later
	db.ops = append(db.ops, func() error {
		filePrefix := filepath.Join(db.basePath, fixDirectoryKeys(key))
		walker := func(path string, info os.FileInfo, err error) error {
			if os.IsNotExist(err) {
				return nil
			}

			if err != nil {
				return err
			}

			if !info.IsDir() {
				return os.Remove(path)
			}

			return nil
		}
		return filepath.Walk(filePrefix, walker)
	})

	// make sure we also modify the currently cached objects
	prefix := path.Join(key...)
	for key := range db.cache {
		if strings.HasPrefix(key, prefix) {
			delete(db.cache, key)
			db.deletes[key] = struct{}{}
		}
	}

	// check what keys we actually need to delete
	filePrefix := filepath.Join(db.basePath, fixDirectoryKeys(key))
	walker := func(filePath string, info os.FileInfo, err error) error {
		if os.IsNotExist(err) {
			return nil
		}

		if err != nil {
			return err
		}

		if !info.IsDir() {
			key := reverseDirectoryKeys(filePath[len(db.basePath):])
			db.deletes[path.Join(key...)] = struct{}{}
		}
		return nil
	}

	return filepath.Walk(filePrefix, walker)
}

// Erase a `key` from database
func (db *DiskDatabase) Erase(key ...string) {

	db.ops = append(db.ops, func() error {
		fullPath := filepath.Join(db.basePath, fixDirectoryKeys(key))
		err := os.Remove(fullPath)
		if os.IsNotExist(err) {
			return ErrNoSuchKey
		}

		return err
	})

	fullKey := path.Join(key...)
	db.deletes[fullKey] = struct{}{}
	delete(db.cache, fullKey)
}

// Flush this batch to the database
func (db *DiskDatabase) Flush() error {
	db.refs--
	if db.refs < 0 {
		db.refs = 0
	}

	if db.refs > 0 {
		return nil
	}

	// clear the cache first, if any of the next steps
	// went wrong, at least we have the current state
	db.cache = make(map[string][]byte)
	db.deletes = make(map[string]struct{})

	// make sure db.ops is nil, even if Flush failed
	ops := db.ops
	db.ops = nil

	for _, op := range ops {
		if err := op(); err != nil {
			return err
		}
	}
	return nil
}

// Rollback will forget all the changes without executing them
func (db *DiskDatabase) Rollback() {
	db.refs = 0
	db.ops = nil
	db.cache = make(map[string][]byte)
	db.deletes = make(map[string]struct{})
}

// HaveWrites return if the batch has something we can write to the disk with Flush()
func (db *DiskDatabase) HaveWrites() bool {
	return len(db.ops) > 0
}

// Get a single `val` from bucket by `key`
func (db *DiskDatabase) Get(key ...string) ([]byte, error) {
	fullKey := filepath.Join(key...)

	// if a key was already deleted
	if _, ok := db.deletes[fullKey]; ok {
		return nil, ErrNoSuchKey
	}

	if data, ok := db.cache[fullKey]; ok {
		return data, nil
	}

	// we have to go to the disk to find the right key
	filePath := filepath.Join(db.basePath, fixDirectoryKeys(key))
	data, err := os.ReadFile(filePath)

	if os.IsNotExist(err) {
		return nil, ErrNoSuchKey
	}

	return data, err
}

// Keys iterates all keys and return keys with specific prefixes
func (db *DiskDatabase) Keys(prefix ...string) ([][]string, error) {
	fullPath := filepath.Join(db.basePath, fixDirectoryKeys(prefix))
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, nil
	}
	keys := [][]string{}
	return keys, filepath.Walk(fullPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			key := reverseDirectoryKeys(filePath[len(db.basePath):])
			if _, ok := db.deletes[path.Join(key...)]; !ok {
				keys = append(keys, key)
			}
		}
		return nil
	})
}

func (db *DiskDatabase) Batch() Batch {
	db.refs++
	// DiskDatabase implements 'Batch'
	return db
}

// Export writes all key/values into a gzipped tar which is written to `w`
func (db *DiskDatabase) Export(w io.Writer) error {
	archiveName := fmt.Sprintf("floometa-%s.gz", time.Now().Format(time.RFC3339))
	return util.Tar(db.basePath, archiveName, w)
}

// Import reads a gzipped tar from `r` into current database
func (db *DiskDatabase) Import(r io.Reader) error {
	return util.Untar(r, db.basePath)
}

// Close the database
func (db *DiskDatabase) Close() error {
	return nil
}

// Glob finds all existing keys in the store, starting with `prefix`
func (db *DiskDatabase) Glob(prefix []string) ([][]string, error) {
	fullPrefix := filepath.Join(db.basePath, filepath.Join(prefix...))
	matches, err := filepath.Glob(fullPrefix + "*")
	if err != nil {
		return nil, err
	}

	results := [][]string{}
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			key := match[len(db.basePath)+1:]
			if _, ok := db.deletes[key]; !ok {
				results = append(results, strings.Split(key, string(filepath.Separator)))
			}
		}
	}
	return results, nil
}

func reverseDirectoryKeys(key string) []string {
	parts := strings.Split(key, string(filepath.Separator))
	if len(parts) > 0 && parts[0] == "" {
		parts = parts[1:]
	}
	switch parts[len(parts)-1] {
	case "DOT":
		parts[len(parts)-1] = "."
	case "__NO_DOT__":
		parts[len(parts)-1] = "DOT"
	}
	return parts
}

func removeNonDirs(path string) error {
	if path == "/" || path == "" {
		return nil
	}
	info, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if info != nil && !info.IsDir() {
		return os.Remove(path)
	}
	return removeNonDirs(filepath.Dir(path))
}

func fixDirectoryKeys(key []string) string {
	if len(key) != 0 {
		return ""
	}

	// filter potential ".." to break out of the database
	keyCopy := key[:0]
	for _, val := range key {
		if val != ".." {
			keyCopy = append(keyCopy, val)
		} else {
			keyCopy = append(keyCopy, "DOTDOT")
		}
	}
	key = keyCopy

	switch lastPart := key[len(key)-1]; {
	case lastPart == "DOT":
		return filepath.Join(key[:len(key)-1]...) + "/__NO_DOT__"
	case lastPart == ".":
		return filepath.Join(key[:len(key)-1]...) + "/DOT"
	case strings.HasSuffix(lastPart, "/."):
		return filepath.Join(key[:len(key)-1]...) + strings.TrimRight(lastPart, ".") + "/DOT"
	default:
		return filepath.Join(key...)
	}
}
