package deploy

import (
	"errors"

	"github.com/go-xiaohei/pugo-static/app/builder"
	"gopkg.in/ini.v1"
)

type (
	Deployer struct {
		tasks []DeployTask
	}
	DeployTask interface {
		New(name string, section *ini.Section) (DeployTask, error) // new instance
		Name() string                                              // task name, ini conf.ini
		Type() string                                              // type string
		IsValid() error                                            // is option valid
		Do(b *builder.Builder, ctx *builder.Context) error         // deploy logic
	}
)

func New(section *ini.Section) (*Deployer, error) {
	items := section.KeysHash()
	if len(items) == 0 {
		return nil, errors.New("please write deploy settings in conf.ini")
	}
	d := &Deployer{}
	for name, item := range items {
		println(name, item)
	}
	return d, nil
}

func (dp *Deployer) Run() error {
	println("run")
	return nil
}

func (dp *Deployer) RunAsync() error {
	println("run async")
	return nil
}
