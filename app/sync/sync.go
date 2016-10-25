package sync

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/helper"
	"gopkg.in/inconshreveable/log15.v2"
)

// Syncer manage sync file
type Syncer struct {
	dir         string
	syncLock    sync.Mutex
	syncedFiles map[string]bool
}

// NewSyncer create a sync object to sync file to dir
func NewSyncer(dir string) *Syncer {
	return &Syncer{
		dir:         dir,
		syncedFiles: make(map[string]bool),
	}
}

// SyncForce force to write file to dst file
func (s *Syncer) SyncForce() error {
	return nil
}

// Sync write to new file to old file if md5 changes
func (s *Syncer) Sync() error {
	return nil
}

// DirOption set options when sync directory
type DirOption struct {
	Filter func(string) bool
	Prefix string
	Ignore []string
}

// SyncDir sync directory files to syncer's directory
func (s *Syncer) SyncDir(dir string, opt *DirOption) error {
	if !com.IsDir(dir) {
		return nil
	}
	var (
		relFile string
		dstFile string
	)
	return filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if opt != nil && opt.Filter != nil {
			if !opt.Filter(p) {
				return nil
			}
		}
		relFile, _ = filepath.Rel(dir, p)
		if opt != nil {
			if len(opt.Ignore) > 0 {
				for _, ignore := range opt.Ignore {
					if strings.HasPrefix(relFile, ignore) {
						return nil
					}
				}
			}
			if opt.Prefix != "" {
				relFile = filepath.Join(opt.Prefix, relFile)
			}
		}
		dstFile = filepath.Join(s.dir, relFile)
		if com.IsFile(dstFile) {
			hash1, _ := helper.Md5File(p)
			hash2, _ := helper.Md5File(dstFile)
			if hash1 == hash2 {
				log15.Debug("Sync|Keep|%s", dstFile)
				s.SetSynced(dstFile)
				return nil
			}
		}
		os.MkdirAll(filepath.Dir(dstFile), os.ModePerm)
		if err := com.Copy(p, dstFile); err != nil {
			return err
		}
		log15.Debug("Sync|Write|%s", dstFile)
		s.SetSynced(dstFile)
		return nil
	})
}

// SetSynced set file as synced file
func (s *Syncer) SetSynced(file string) {
	file = filepath.ToSlash(file)
	s.syncLock.Lock()
	s.syncedFiles[file] = true
	s.syncLock.Unlock()
}

// Clear clean all non-synced file in s.dir
func (s *Syncer) Clear() error {
	return filepath.Walk(s.dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		p = filepath.ToSlash(p)
		if s.syncedFiles[p] {
			return nil
		}
		log15.Debug("Sync|Del|%s", p)
		return os.Remove(p)
	})
}
