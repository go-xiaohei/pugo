package migrate

import (
	"strings"

	"github.com/go-xiaohei/pugo/app/builder"
	"gopkg.in/inconshreveable/log15.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Init init migration handler to builder
func Init() {
	builder.Before(Handle)
}

// Handle detect ctx to get correct migrate Task and run it
func Handle(ctx *builder.Context) {
	var (
		task Task
		err  error
	)
	for _, m := range manager.tasks {
		task, err = m.Detect(ctx)
		if err != nil {
			log15.Crit("Migrate|Fail|%s", err.Error())
			return
		}
		if task != nil {
			log15.Info("Migrate|Detect|%s", task.Name())
			files, err := task.Action(ctx)
			if err != nil {
				log15.Error("Migrate|%s|%s", task.Name(), err.Error())
				return
			}
			for file, b := range files {
				file = filepath.Join(ctx.SrcDir(), file)
				log15.Debug("Migrate|Write|%s", file)
				if b == nil {
					os.MkdirAll(file, os.ModePerm)
					continue
				}
				os.MkdirAll(filepath.Dir(file), os.ModePerm)
				if err = ioutil.WriteFile(file, b.Bytes(), os.ModePerm); err != nil {
					log15.Error("Migrate|%s|%s", task.Name(), err.Error())
					return
				}
			}
		}
	}
	if isMigrateTo(ctx.From) {
		log15.Crit("Migrate|Unknown|%s", ctx.From)
	}
}

func isMigrateTo(to string) bool {
	if !strings.Contains(to, "://") {
		return false
	}
	if strings.HasPrefix(to, "dir://") {
		return false
	}
	return true
}
