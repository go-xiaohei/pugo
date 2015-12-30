package command

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/builder"
	"github.com/go-xiaohei/pugo/app/deploy"
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
		Action: buildSite(opt, false),
		Before: setDebugMode,
	}
}

// build site function
func buildSite(opt *builder.BuildOption, mustWatch bool) func(ctx *cli.Context) {
	return func(ctx *cli.Context) {
		// ctrl+C capture
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		if theme := ctx.String("theme"); theme != "" {
			opt.Theme = theme
		}
		if opt.TargetDir == "" {
			// set target dir
			if targetDir := ctx.String("dest"); targetDir != "" {
				opt.TargetDir = targetDir
			}
			if opt.TargetDir == "template" || opt.TargetDir == "source" {
				log15.Crit("Builder.Fail", "error", "destination directory should not be 'template' or 'source'")
			}
		}

		b := builder.New(opt)
		if b.Error != nil {
			log15.Crit("Builder.Fail", "error", b.Error.Error())
		}

		// detect deploy callback
		way, err := deploy.Detect(ctx)
		if err != nil {
			log15.Crit("Deploy.Fail", "error", err.Error())
		}
		if way != nil {
			opt.TargetDir = way.Dir()
			opt.After(func(b *builder.Builder, ctx *builder.Context) error {
				t := time.Now()
				if err := way.Do(b, ctx); err != nil {
					return err
				}
				log15.Info("Deploy.Finish", "duration", time.Since(t))
				return nil
			})
		}

		// make directory
		log15.Info("Dest." + opt.TargetDir)
		if com.IsDir(opt.TargetDir) {
			log15.Warn("Dest." + opt.TargetDir + ".Existed")
		}

		// auto watching
		b.Build(opt.TargetDir)
		if err := b.Context().Error; err != nil {
			log15.Crit("Build.Fail", "error", err.Error())
		}

		if ctx.Bool("watch") || mustWatch {
			b.Watch(opt.TargetDir)
			<-signalChan
		}

		log15.Info("Build.Close")
	}
}
