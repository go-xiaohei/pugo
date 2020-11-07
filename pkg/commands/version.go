package commands

import (
	"fmt"
	"os"
	"pugo/pkg/vars"
	"runtime"

	"github.com/urfave/cli/v2"
)

var (
	// Version is command of 'version'
	Version = &cli.Command{
		Name:  "version",
		Usage: "print PuGo Version info",
		Action: func(ctx *cli.Context) error {
			execFile, _ := os.Executable()
			info, _ := os.Stat(execFile)
			fmt.Printf("%v version:\t %s\n", ctx.App.Name, ctx.App.Version)
			fmt.Printf("Go version:\t %s\n", runtime.Version())
			fmt.Printf("Binary size:\t %.2fMB\n", float64(info.Size())/1024/1024)
			fmt.Printf("Compiled time:\t %s\n", ctx.App.Compiled.Format("2006/01/02 15:04:05"))
			if vars.Commit != "" {
				fmt.Printf("Commit:\t %s\n", vars.Commit)
			}
			fmt.Printf("OS/Arch:\t %s\n", runtime.GOOS+"/"+runtime.GOARCH)
			return nil
		},
	}
)
