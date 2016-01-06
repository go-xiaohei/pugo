package command

import (
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/migrate"
	"gopkg.in/inconshreveable/log15.v2"
)

// Migrate is a command to migrate from other content system
func Migrate() cli.Command {
	return cli.Command{
		Name:     "migrate",
		Usage:    "migrate content from other system",
		HideHelp: true,
		Flags: []cli.Flag{
			srcFlag,
			destFlag,
			debugFlag,
		},
		Action: migrateSite(),
		Before: setDebugMode,
	}
}

func migrateSite() func(ctx *cli.Context) {
	return func(ctx *cli.Context) {
		task, err := migrate.Detect(ctx)
		if err != nil {
			log15.Crit("Migrate.Fail", "error", err.Error())
		}
		if task == nil {
			log15.Crit("Migrate.Fail", "error", migrate.ErrMigrateUnknown.Error())
		}
		if err = task.Do(); err != nil {
			log15.Crit("Migrate.Fail", "error", err.Error())
		}
		log15.Info("Migrate.Done.[" + task.Type() + "]")
	}
}
