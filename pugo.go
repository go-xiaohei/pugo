package main

import (
	"time"

	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/command"
)

const (
	// Name is app name
	Name = "PuGo"
	// Desc is app usage
	Desc = "A Fast Static Site Generator"
	// Version is app version number
	Version = "0.10.0 (beta)"
)

func main() {
	app := cli.NewApp()
	app.Name = Name
	app.Usage = Desc
	app.Version = Version
	app.Compiled = time.Now()
	app.Before = command.Before
	app.Commands = []cli.Command{
		command.Build,
	}
	app.RunAndExitOnError()
}
