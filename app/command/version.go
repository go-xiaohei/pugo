package command

import (
	"fmt"
	"runtime"

	"github.com/go-xiaohei/pugo/app/vars"
	"github.com/urfave/cli"
)

var (
	// Version is command of 'version'
	Version = cli.Command{
		Name:  "version",
		Usage: "print PuGo Version",
		Action: func(ctx *cli.Context) error {
			fmt.Printf("%v version :\t %s\n", ctx.App.Name, ctx.App.Version)
			fmt.Printf("Go version :\t %s\n", runtime.Version())
			fmt.Printf("Compiled time :\t %s\n", ctx.App.Compiled.Format("2006/01/02 15:04"))

			if vars.Commit != "" {
				fmt.Printf("Commit :\t %s\n", vars.Commit)
			}
			return nil
		},
	}
)
