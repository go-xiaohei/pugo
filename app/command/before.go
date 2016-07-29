package command

import (
	"os"

	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/urfave/cli"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/inconshreveable/log15.v2/ext"
)

// Before set before handler when start run cli.App
func Before(ctx *cli.Context) error {
	lv := log15.LvlInfo
	if ctx.Bool("debug") {
		lv = log15.LvlDebug
	}
	log15.Root().SetHandler(log15.LvlFilterHandler(lv, ext.FatalHandler(log15.StreamHandler(os.Stderr, helper.LogfmtFormat()))))
	return nil
}
