package main

import (
	"time"

	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/command"
	"github.com/go-xiaohei/pugo/app/vars"
)

func main() {
	app := cli.NewApp()
	app.Name = vars.Name
	app.Usage = vars.Desc
	app.Version = vars.Version
	app.Compiled = time.Now()
	app.Commands = []cli.Command{
		command.Build,
	}
	app.RunAndExitOnError()
}
