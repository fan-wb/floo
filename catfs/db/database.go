package db

import (
	"errors"
	"io"
)

type Batch interface {
	// Put sets `val` at `key`
	Put(val []byte, key ...string)

	// Clear all contents below and including `key`
	Clear(key ...string) error

	// Erase a `key` from database
	Erase(key ...string)

	// Flush the batch to database
	// this is when all changes are written to the disk
	Flush() error

	// Rollback will forget all the changes without executing them
	Rollback()

	// HaveWrites return if the batch has something we can
	// write to the disk with Flush()
	HaveWrites() bool
}

type Database interface {
	// Get retrieves `key` out of the bucket
	Get(key ...string) ([]byte, error)

	// Keys iterates all keys in the database and return in lexical order
	Keys(prefix ...string) ([][]string, error)

	// Batch returns new Batch object
	Batch() Batch

	// Export backups all database content to `w`
	// in an implementation specific format that can be read by Import
	Export(w io.Writer) error

	// Import reads a previously exported db dump by Export()
	// Existing keys might be overwritten if the dump also contains them
	Import(r io.Reader) error

	// Close a database
	Close() error

	// Glob finds all existing keys in the store starts with `prefix`
	Glob(prefix string) ([][]string, error)
}

var (
	ErrNoSuchKey = errors.New("The key does not exist")
)
