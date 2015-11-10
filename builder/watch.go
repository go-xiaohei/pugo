package builder

import (
	"gopkg.in/fsnotify.v1"
	"gopkg.in/inconshreveable/log15.v2"
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
				if path.Ext(event.Name) == ".md" {
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

	watcher.Add(b.srcDir)
}
