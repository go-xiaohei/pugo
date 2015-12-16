package deploy

import (
	"time"

	"github.com/go-xiaohei/pugo-static/app/builder"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/ini.v1"
	"sync"
)

type (
	// Deployer contains tasks and registered task types
	Deployer struct {
		tasks      []DeployTask
		registered map[string]DeployTask
	}
	// DeployTask defines the methods of a deploy task
	DeployTask interface {
		New(name string, section *ini.Section) (DeployTask, error) // new instance
		Name() string                                              // task name, ini conf.ini
		Type() string                                              // type string
		IsValid() error                                            // is option valid
		Do(b *builder.Builder, ctx *builder.Context) error         // deploy logic
	}
)

// New Deployer with conf ini file
func New() *Deployer {
	return &Deployer{
		registered: map[string]DeployTask{
			TYPE_GIT:  new(GitTask),
			TYPE_FTP:  new(FtpTask),
			TYPE_SFTP: new(SftpTask),
		},
	}
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
				log15.Error("Deploy.Task.["+task.Name()+"]", "error", err.Error(), "duration", time.Since(t))
			} else {
				log15.Info("Deploy.Task.["+task.Name()+"]", "duration", time.Since(t))
			}
		}(task)
	}
	wg.Wait()
	log15.Info("Deploy.Finish", "duration", time.Since(t))
	return nil
}
