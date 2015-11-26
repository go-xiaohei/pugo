package command

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo-static/builder"
	"gopkg.in/inconshreveable/log15.v2"
)

func Build(opt *builder.BuildOption) cli.Command {
	return cli.Command{
		Name:     "build",
		Usage:    "build static files and watch updating",
		HideHelp: true,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "target",
				Usage: "build static files to target directory",
				Value: "dest",
			},
			cli.BoolFlag{
				Name:  "watch",
				Usage: "watch source changes to rebuild",
			},
			cli.StringFlag{
				Name:  "theme",
				Value: "default",
				Usage: "set theme to build",
			},
		},
		Action: buildSite(opt),
	}
}

func buildSite(opt *builder.BuildOption) func(ctx *cli.Context) {
	return func(ctx *cli.Context) {
		// ctrl+C capture
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		opt.Theme = ctx.String("theme")
		opt.IsDebug = true
		opt.IsWatchTemplate = true
		opt.IsCopyAssets = true

		log15.Debug("Dir.Source." + opt.SrcDir)
		log15.Debug("Dir.Template." + opt.TplDir)

		b := builder.New(opt)
		if b.Error != nil {
			log15.Crit("Build.Fail", "error", b.Error.Error())
		}

		targetDir := ctx.String("target")
		log15.Info("Dir.Target." + targetDir)
		if com.IsDir(targetDir) {
			log15.Warn("Dir.Target." + targetDir + ".Existed")
		}
		b.Build(targetDir)
		if err := b.Report().Error; err != nil {
			log15.Crit("Build.Fail", "error", err.Error())
		}
		if ctx.Bool("watch") {
			log15.Info("Build.Watch")
			b.Watch(targetDir)
			<-signalChan
		}

		log15.Info("Build.Close")
	}
}
