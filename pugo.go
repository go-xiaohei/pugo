package main

import (
	"time"

	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/command"
	"github.com/go-xiaohei/pugo/app/vars"
)

//go:generate go-bindata -o=app/asset/asset.go -pkg=asset source/... theme/... doc/source/... doc/theme/...
//go:generate gofmt -w -s .
//go:generate git log --format=%h -n 1

func main() {
	app := cli.NewApp()
	app.Name = vars.Name
	app.Usage = vars.Desc
	app.Version = vars.Version
	app.Compiled = time.Now()
	app.Commands = []cli.Command{
		command.Build,
		command.Server,
		command.New,
		command.Doc,
		command.Deploy,
		command.Version,
	}
	app.HideVersion = true
	app.RunAndExitOnError()
}
