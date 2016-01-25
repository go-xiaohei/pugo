package deploy

import (
	"strings"

	"github.com/go-xiaohei/pugo/app/builder"
	"gopkg.in/inconshreveable/log15.v2"
)

// Init init deployment handler to builder
func Init() {
	builder.Before(Detect)
	builder.After(Action)
}

// Detect detect ctx to get correct deploy Task
func Detect(ctx *builder.Context) {
	var (
		task Task
		err  error
	)
	for _, m := range manager.tasks {
		task, err = m.Detect(ctx)
		if err != nil {
			log15.Crit("Deploy|Fail|%s", err.Error())
			return
		}
		if isDeployTo(ctx.To) && task == nil {
			log15.Warn("Deploy|Unknown|%s", ctx.To)
		}
		if task != nil {
			log15.Info("Deploy|Detect|%s", task.Name())
			manager.Set(task)
		}
	}
}

// Action use detected Task to run Task.Action
func Action(ctx *builder.Context) {
	task := manager.Get()
	if task == nil {
		return
	}
	log15.Info("Deploy|Action|%s", task.Name())
	if err := task.Action(ctx); err != nil {
		log15.Error("Deploy|Fail|%s", err.Error())
	}
}

func isDeployTo(to string) bool {
	if !strings.Contains(to, "://") {
		return false
	}
	if strings.HasPrefix(to, "dir://") {
		return false
	}
	return true
}
