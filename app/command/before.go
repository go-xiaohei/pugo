package command

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/helper"
	"gopkg.in/inconshreveable/log15.v2"
)

// Before set before handler when start run cli.App
func Before(ctx *cli.Context) error {
	log15.Root().SetHandler(log15.StreamHandler(os.Stderr, helper.LogfmtFormat()))
	return nil
}
