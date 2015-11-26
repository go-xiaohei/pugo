package command

import (
	"github.com/codegangsta/cli"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/inconshreveable/log15.v2/ext"
)

var (
	debugFlag = cli.BoolFlag{
		Name:  "debug",
		Usage: "print debug logs",
	}
	themeFlag = cli.StringFlag{
		Name:  "theme",
		Value: "default",
		Usage: "set theme to render",
	}
	destFlag = cli.StringFlag{
		Name:  "dest",
		Value: "dest",
		Usage: "set compiling to directory",
	}
)

func setDebugMode(ctx *cli.Context) error {
	if ctx.Bool("debug") {
		log15.Root().SetHandler(log15.LvlFilterHandler(log15.LvlDebug, ext.FatalHandler(log15.StderrHandler)))
	}
	return nil
}
