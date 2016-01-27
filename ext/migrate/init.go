package migrate

import (
	"github.com/go-xiaohei/pugo/app/builder"
	"gopkg.in/inconshreveable/log15.v2"
	"strings"
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
			return task.Action(ctx)
		}
		if isMigrateTo(ctx.From) && task == nil {
			log15.Warn("Migrate|Unknown|%s", ctx.From)
		}
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
