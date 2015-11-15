package builder

import (
	"gopkg.in/fsnotify.v1"
	"gopkg.in/inconshreveable/log15.v2"
	"io/ioutil"
	"path"
)

func (b *Builder) Watch(dstDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log15.Crit("Build.Watch", "error", err.Error())
		return
	}
	log15.Info("Build.Watch")

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				ext := path.Ext(event.Name)
				if ext == ".md" || ext == ".html" {
					if event.Op != fsnotify.Chmod && !b.IsBuilding() {
						log15.Info("Build.Watch.Rebuild", "change", event.String())
						b.Build(dstDir)
					}
				}
			case err := <-watcher.Errors:
				log15.Error("Build.Watch", "error", err.Error())
			}
		}
	}()

	watchDir(watcher, b.srcDir)
	if b.opt.IsWatchTemplate {
		watchDir(watcher, b.tplDir)
	}
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
