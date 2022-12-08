package fuse

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"floo/util"
)

// MountOptions defines all possible knobs you can turn for a mount.
// The zero value are the default options.
type MountOptions struct {
	// ReadOnly makes the mount not modifyable
	ReadOnly bool
	// Root determines what the root directory is.
	Root string
	// Offline tells the mount to error out on files that would need
	// to be fetched from far.
	Offline bool
}

type Mount struct {
	Dir string

	filesys *Filesystem
	closed  bool
	done    chan util.Empty
	errors  chan error
	conn    *fuse.Conn
	server  *fs.Server
	options MountOptions
	// notifier Notifier
	// fs       *catfs.FS
}
