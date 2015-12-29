package builder

import (
	"io/ioutil"
	"path"
	"time"

	"gopkg.in/fsnotify.v1"
	"gopkg.in/inconshreveable/log15.v2"
	"path/filepath"
)

var (
	watchingExt       = []string{".md", ".ini", ".html", ".css", ".js"}
	watchScheduleTime int64
)

// watch source dir changes and build to destination directory if updated
func (b *Builder) Watch(dstDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log15.Crit("Watch.Fail", "error", err.Error())
		return
	}
	log15.Info("Watch.Start")
	b.isWatching = true

	// use a ticker to trigger build
	go func() {
		c := time.Tick(1 * time.Second)
		for {
			t := <-c
			if watchScheduleTime > 0 && t.UnixNano() > watchScheduleTime {
				b.Build(dstDir)
				watchScheduleTime = 0
			}
		}
	}()

	// catch fsnotify events
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				ext := path.Ext(event.Name)
				for _, e := range watchingExt {
					if e == ext {
						if event.Op != fsnotify.Chmod && !b.IsBuilding() {
							log15.Info("Watch.Rebuild", "change", filepath.ToSlash(event.String()))
							watchScheduleTime = time.Now().Add(time.Second).UnixNano()
						}
						break
					}
				}
			case err := <-watcher.Errors:
				log15.Error("Watch.Errors", "error", err.Error())
			}
		}
	}()

	watchDir(watcher, b.opt.SrcDir)
	watchDir(watcher, b.opt.TplDir)
}

// watch sub directory
func watchDir(watcher *fsnotify.Watcher, srcDir string) {
	watcher.Add(srcDir)
	dir, _ := ioutil.ReadDir(srcDir)
	for _, d := range dir {
		if d.IsDir() {
			watchDir(watcher, path.Join(srcDir, d.Name()))
		}
	}
}
