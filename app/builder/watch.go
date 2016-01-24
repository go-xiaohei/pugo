package builder

import (
	"io/ioutil"
	"path"
	"time"

	"gopkg.in/fsnotify.v1"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	// watchingExt sets the suffix that watching to
	watchingExt = []string{".md", ".toml", ".html", ".css", ".js", ".jpg", ".png", ".gif"}
	// watchScheduleTime sets watching timer duration
	watchScheduleTime int64
)

// Watch watch changes
func Watch(ctx *Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log15.Crit("Watch|Fail|%s", err.Error())
		return
	}
	log15.Info("Watch|Start")

	if ctx.srcDir == "" || ctx.dstDir == "" || ctx.Theme == nil {
		log15.Crit("Watch|Need build once then watch changes")
	}

	// use a ticker to trigger build
	go func() {
		c := time.Tick(1 * time.Second)
		for {
			t := <-c
			if watchScheduleTime > 0 && t.UnixNano() > watchScheduleTime {
				ctx.Again()
				Build(ctx)
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
						if event.Op != fsnotify.Chmod {
							log15.Info("Watch|Rebuild|%s", event.String())
							watchScheduleTime = time.Now().Add(time.Second).UnixNano()
						}
						break
					}
				}
			case err := <-watcher.Errors:
				log15.Warn("Watch|Error|%s", err.Error())
			}
		}
	}()

	watchDir(watcher, ctx.srcDir)
	watchDir(watcher, ctx.Theme.Dir())

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
