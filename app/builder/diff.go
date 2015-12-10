package builder

import (
	"path/filepath"
	"time"
)

const (
	DIFF_ADD    = 1
	DIFF_UPDATE = 2
	DIFF_KEEP   = 3 // keep old file, do not change it

	DIFF_REMOVE = 9
)

type (
	Diff struct {
		files map[string]*DiffEntry
	}
	DiffEntry struct {
		behavior int
		t        time.Time
	}
)

func newDiff() *Diff {
	return &Diff{
		files: make(map[string]*DiffEntry),
	}
}

func (d *Diff) Add(file string, behavior int) {
	file = filepath.ToSlash(file)
	d.files[file] = &DiffEntry{behavior: behavior, t: time.Now()}
}

func (d *Diff) Exist(file string) bool {
	file = filepath.ToSlash(file)
	_, ok := d.files[file]
	return ok
}
