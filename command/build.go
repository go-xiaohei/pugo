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

// build Command, need BuildOption
func Build(opt *builder.BuildOption) cli.Command {
	return cli.Command{
		Name:     "build",
		Usage:    "build static files and watch updating",
		HideHelp: true,
		Flags: []cli.Flag{
			destFlag,
			themeFlag,
			watchFlag,
			debugFlag,
		},
		Action: buildSite(opt),
		Before: setDebugMode,
	}
}

// build site function
func buildSite(opt *builder.BuildOption) func(ctx *cli.Context) {
	return func(ctx *cli.Context) {
		// ctrl+C capture
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		opt.Theme = ctx.String("theme")

		b := builder.New(opt)
		if b.Error != nil {
			log15.Crit("Builder.Fail", "error", b.Error.Error())
		}

		targetDir := ctx.String("dest")
		log15.Info("Dest." + targetDir)
		if com.IsDir(targetDir) {
			log15.Warn("Dest." + targetDir + ".Existed")
		}

		// auto watching
		b.Build(targetDir)
		if err := b.Context().Error; err != nil {
			log15.Crit("Build.Fail", "error", err.Error())
		}

		if ctx.Bool("watch") {
			b.Watch(targetDir)
			<-signalChan
		}

		log15.Info("Build.Close")
	}
}
