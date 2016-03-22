package command

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/vars"
)

var (
	// Version is command of 'version'
	Version = cli.Command{
		Name:  "version",
		Usage: "print PuGo Version",
		Action: func(ctx *cli.Context) {
			if vars.Commit == "" {
				fmt.Printf("%v version %v ~ %s\n", ctx.App.Name, ctx.App.Version, ctx.App.Compiled.Format("2006/01/02 15:04"))
				return
			}
			fmt.Printf("%v version %v ~ %s\n", ctx.App.Name, ctx.App.Version, vars.Commit)
		},
	}
)
