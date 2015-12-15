package deploy

import (
	"errors"

	"github.com/go-xiaohei/pugo-static/app/builder"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/ini.v1"
	"time"
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
func New(file *ini.File) (*Deployer, error) {
	items := file.Section("deploy").KeysHash()
	if len(items) == 0 {
		return nil, errors.New("please write deploy settings in conf.ini")
	}
	d := &Deployer{
		registered: map[string]DeployTask{
			TYPE_GIT:  new(GitTask),
			TYPE_FTP:  new(FtpTask),
			TYPE_SFTP: new(SftpTask),
		},
	}
	for _, name := range items {
		s := file.Section("deploy." + name)
		typeName := s.Key("type").String()
		if _, ok := d.registered[typeName]; !ok {
			return nil, errors.New("unsupport deploy type : " + typeName)
		}
		task, err := d.registered[typeName].New(name, s)
		if err != nil {
			return nil, err
		}
		d.tasks = append(d.tasks, task)
		log15.Debug("Deploy.AddTask.[" + name + "]")
	}
	return d, nil
}

// run deployer tasks,
// if error, return error and stop
func (dp *Deployer) Run(b *builder.Builder, ctx *builder.Context) error {
	for _, task := range dp.tasks {
		t := time.Now()
		if err := task.Do(b, ctx); err != nil {
			log15.Error("Deploy.Task.["+task.Name()+"]", "error", err.Error(), "duration", time.Since(t))
			return err
		}
		log15.Info("Deploy.Task.["+task.Name()+"]", "duration", time.Since(t))
	}
	return nil
}

// run deployer tasks in goroutine
// if error, just log
func (dp *Deployer) RunAsync(b *builder.Builder, ctx *builder.Context) error {
	for _, task := range dp.tasks {
		go func(task DeployTask) {
			t := time.Now()
			if err := task.Do(b, ctx); err != nil {
				log15.Error("Deploy.Task.["+task.Name()+"]", "error", err.Error(), "duration", time.Since(t))
				return
			}
			log15.Info("Deploy.Task.["+task.Name()+"]", "duration", time.Since(t))
		}(task)
	}
	return nil
}
