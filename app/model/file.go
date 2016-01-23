package model

import (
	"fmt"
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
)

type (
	// File describe a generated or compiled file
	File struct {
		URL     string
		ModTime time.Time
		Size    int64
		Type    string
		Hash    string
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

// Add add file
func (fs *Files) Add(url string, size int64, modTime time.Time, t string) {
	f := &File{
		URL:     url,
		Size:    size,
		ModTime: modTime,
		Type:    t,
	}
	hash, _ := helper.Md5File(url)
	f.Hash = hash
	fs.files[url] = f
}

// Print print files
func (fs *Files) Print() {
	for _, f := range fs.files {
		fmt.Printf("%s %d %s %s@%s\n", f.URL, f.Size, f.ModTime, f.Hash, f.Type)
	}
}
