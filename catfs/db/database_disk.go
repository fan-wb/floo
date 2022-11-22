package db

import (
	"os"
	"path/filepath"
	"strings"
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
