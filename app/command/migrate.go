package command

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/migrate"
	"gopkg.in/inconshreveable/log15.v2"
)

// Migrate is a command to migrate from other content system
func Migrate(toDir string) cli.Command {
	return cli.Command{
		Name:     "migrate",
		Usage:    "migrate content from other system",
		HideHelp: true,
		Flags: []cli.Flag{
			srcFlag,
			toFlag,
			debugFlag,
		},
		Action: migrateSite(toDir),
		Before: setDebugMode,
	}
}

func migrateSite(toDir string) func(ctx *cli.Context) {
	migrate.OutputDirectory = toDir
	return func(ctx *cli.Context) {
		t := time.Now()
		if dest := ctx.String("to"); dest != "" {
			os.MkdirAll(dest, os.ModePerm)
			migrate.OutputDirectory = dest
		}
		log15.Info("Migrate.To", "dir", migrate.OutputDirectory)
		task, err := migrate.Detect(ctx)
		if err != nil {
			log15.Crit("Migrate.Fail", "error", err.Error())
		}
		if task == nil {
			log15.Crit("Migrate.Fail", "error", migrate.ErrMigrateUnknown.Error())
		}
		files, err := task.Do()
		if err != nil {
			log15.Crit("Migrate.Fail", "error", err.Error())
		}
		for filename, b := range files {
			file := path.Join(migrate.OutputDirectory, filename)
			if com.IsFile(file) {
				log15.Warn("Migrate.Conflict", "file", file)
			}
			os.MkdirAll(path.Dir(file), os.ModePerm)
			if b != nil {
				ioutil.WriteFile(file, b.Bytes(), os.ModePerm)
			}
		}
		log15.Info("Migrate.Done.["+task.Type()+"]", "duration", time.Since(t))
	}
}
