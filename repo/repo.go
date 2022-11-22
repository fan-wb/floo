package repo

import "sync"

type Repository struct {
	mu sync.Mutex
	// fsMap map[string]*catfs.FS
}
