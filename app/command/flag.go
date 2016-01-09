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
	srcFlag = cli.StringFlag{
		Name:  "src",
		Value: "",
		Usage: "migrate from source",
	}
	addrFlag = cli.StringFlag{
		Name:  "addr",
		Value: "0.0.0.0:9899",
		Usage: "set http server address",
	}
	watchFlag = cli.BoolFlag{
		Name:  "watch",
		Usage: "watch changes and auto-rebuild",
	}
	toFlag = cli.StringFlag{
		Name:  "to",
		Usage: "output to directory",
	}
)

func setDebugMode(ctx *cli.Context) error {
	if ctx.Bool("debug") {
		log15.Root().SetHandler(log15.LvlFilterHandler(log15.LvlDebug, ext.FatalHandler(log15.StderrHandler)))
	}
	return nil
}
