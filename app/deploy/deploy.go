package deploy

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-xiaohei/pugo-static/app/builder"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	registeredDeployWay      map[string]DeployTask
	ErrDeployConfFormatError = errors.New("deploy format need be type:conf_string")
)

func init() {
	registeredDeployWay = map[string]DeployTask{
		TYPE_GIT: new(GitTask),
	}
}

type (
	// Deployer contains tasks and registered task types
	Deployer struct {
		tasks []DeployTask
	}
	// DeployTask defines the methods of a deploy task
	DeployTask interface {
		New(config string) (DeployTask, error)             // new instance
		Name() string                                      // task name, identifier
		Do(b *builder.Builder, ctx *builder.Context) error // deploy logic
	}
)

// Add new deploy task with conf string
// if parsed conf error, show error
func (dp *Deployer) Add(conf string) error {
	confData := strings.Split(conf, ":")
	if len(confData) < 2 {
		return ErrDeployConfFormatError
	}
	task, ok := registeredDeployWay[confData[0]]
	if !ok {
		return fmt.Errorf("deploy method '%s' is unsupported", confData[0])
	}
	task, err := task.New(strings.TrimLeft(conf, confData[0]+":"))
	if err != nil {
		return err
	}
	dp.tasks = append(dp.tasks, task)
	return nil
}

// run deployer tasks in goroutine
// if error, just log
func (dp *Deployer) Run(b *builder.Builder, ctx *builder.Context) error {
	if len(dp.tasks) == 0 {
		log15.Warn("Deploy.NoTask")
		return nil
	}
	log15.Debug("Deploy.Start")
	var (
		wg sync.WaitGroup
		t  = time.Now()
	)
	wg.Add(len(dp.tasks))
	for _, task := range dp.tasks {
		go func(task DeployTask) {
			defer wg.Done()
			t := time.Now()
			if err := task.Do(b, ctx); err != nil {
				log15.Error("Deploy.["+task.Name()+"]", "error", err.Error(), "duration", time.Since(t))
			} else {
				log15.Info("Deploy.["+task.Name()+"]", "duration", time.Since(t))
			}
		}(task)
	}
	wg.Wait()
	log15.Info("Deploy.Finish", "duration", time.Since(t))
	return nil
}
