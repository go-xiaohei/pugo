package model

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-xiaohei/pugo/app/helper"
)

const (
	// FileCompiled means the file is from compiled data
	FileCompiled = "compiled"
	// FileStatic means it's copied from static directory
	FileStatic = "static"
	// FileMedia means it's copied from media directory
	FileMedia = "media"

	// OpCompiled means this file is after compiled
	OpCompiled = "compiled"
	// OpCopy mean the file is from copying operation
	OpCopy = "copy"
	// OpKeep means the file is keep, no operation
	OpKeep = "Keep"
	// OpRemove means the files is removed in this process
	OpRemove = "remove"
)

type (
	// File describe a generated or compiled or operated file
	File struct {
		URL     string
		ModTime time.Time
		Size    int64
		Type    string
		Hash    string
		Op      string
	}
	// Files record all relative files
	Files struct {
		files map[string]*File
	}
)

// NewFiles create files group
func NewFiles() *Files {
	return &Files{
		files: make(map[string]*File),
	}
}

func (fs *Files) Get(url string) *File {
	url = filepath.ToSlash(url)
	return fs.files[url]
}

// Add add file
func (fs *Files) Add(url string, size int64, modTime time.Time, t string, op string) {
	url = filepath.ToSlash(url)
	f := &File{
		URL:     url,
		Size:    size,
		ModTime: modTime,
		Type:    t,
		Op:      op,
	}
	if op != OpRemove {
		hash, _ := helper.Md5File(url)
		f.Hash = hash
	}
	fs.files[url] = f
}

// Exist check file existing in operated files
func (fs *Files) Exist(file string) bool {
	for _, f := range fs.files {
		if filepath.ToSlash(f.URL) == filepath.ToSlash(file) {
			return f.Op != OpRemove
		}
	}
	return false
}

// All return all files in Files
func (fs *Files) All() map[string]*File {
	return fs.files
}

// Print print files
func (fs *Files) Print() {
	for _, f := range fs.files {
		fmt.Printf("%s %d %s %s@%s\n", f.URL, f.Size, f.ModTime, f.Hash, f.Type)
	}
}
