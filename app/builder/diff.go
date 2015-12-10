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
		Behavior int
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
	d.files[file] = &DiffEntry{Behavior: behavior, t: time.Now()}
}

func (d *Diff) Exist(file string) bool {
	file = filepath.ToSlash(file)
	_, ok := d.files[file]
	return ok
}

func (d *Diff) Walk(fn func(string, *DiffEntry) error) error {
	var err error
	for name, f := range d.files {
		if err = fn(name, f); err != nil {
			return err
		}
	}
	return nil
}
