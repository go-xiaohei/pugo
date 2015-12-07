package builder

import (
	"io/ioutil"
	"path"

	"gopkg.in/fsnotify.v1"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	watchingExt = []string{".md", ".html", ".css", ".js"}
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

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				ext := path.Ext(event.Name)
				for _, e := range watchingExt {
					if e == ext {
						if event.Op != fsnotify.Chmod && !b.IsBuilding() {
							log15.Info("Watch.Rebuild", "change", event.String())
							b.Build(dstDir)
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

func watchDir(watcher *fsnotify.Watcher, srcDir string) {
	watcher.Add(srcDir)
	dir, _ := ioutil.ReadDir(srcDir)
	for _, d := range dir {
		if d.IsDir() {
			watchDir(watcher, path.Join(srcDir, d.Name()))
		}
	}
}
