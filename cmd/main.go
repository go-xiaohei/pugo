package main

import (
	"os"
	"time"

	"pugo/pkg/commands"
	"pugo/pkg/vars"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        vars.Name,
		Usage:       vars.Desc,
		Version:     vars.Version,
		Compiled:    time.Now(),
		HideVersion: true,
	}
	app.Commands = []*cli.Command{
		commands.Version,
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
