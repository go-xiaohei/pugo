package builder

import (
	"path/filepath"
	"time"
)

const (
	// DiffAdd add new file
	DiffAdd = 1
	// DiffUpdate update file, modified
	DiffUpdate = 2
	// DiffKeep keep old file, do not change it
	DiffKeep = 3
	// DiffRemove remove file
	DiffRemove = 9
)

type (
	// Diff saves diff changes in once context
	Diff struct {
		files map[string]*Entry
	}
	// Entry contains each diff change for a file
	Entry struct {
		Behavior int
		Time     time.Time
	}
)

func newDiff() *Diff {
	return &Diff{
		files: make(map[string]*Entry),
	}
}

// Add adds diff file, behavior int and modification time
func (d *Diff) Add(file string, behavior int, t time.Time) {
	file = filepath.ToSlash(file)
	d.files[file] = &Entry{Behavior: behavior, Time: t}
}

// Exist checks file existing
func (d *Diff) Exist(file string) bool {
	file = filepath.ToSlash(file)
	_, ok := d.files[file]
	return ok
}

// Walk runs all diff Entry in this diff task
func (d *Diff) Walk(fn func(string, *Entry) error) error {
	var err error
	for name, f := range d.files {
		if err = fn(name, f); err != nil {
			return err
		}
	}
	return nil
}
