package command

import (
	"errors"
	"time"

	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo-static/app/builder"
	"github.com/go-xiaohei/pugo-static/app/deploy"
	"gopkg.in/inconshreveable/log15.v2"
)

func Deploy(opt *builder.BuildOption) cli.Command {
	return cli.Command{
		Name:     "deploy",
		Usage:    "deploy site to other platform",
		HideHelp: true,
		Flags: []cli.Flag{
			destFlag,
			themeFlag,
			debugFlag,
			watchFlag,
		},
		Action: deploySite(opt),
		Before: setDebugMode,
	}
}

func deploySite(opt *builder.BuildOption) func(ctx *cli.Context) {
	// real deploy action, in builder hook
	afterFunc := func(dopt *deploy.Option) builder.BuildHook {
		return func(b *builder.Builder, ctx *builder.Context) error {

			// do git deployment
			for name, gitOpt := range dopt.GitOptions {
				fn := func(name string, gitOpt deploy.GitOption) {
					log15.Info("Deploy.Git.[" + name + "].Start")
					t := time.Now()
					if err := deploy.Git(gitOpt, ctx); err != nil {
						log15.Error("Deploy.Git.["+name+"]", "error", err, "duration", time.Since(t))
						return
					}
					log15.Info("Deploy.Git.["+name+"].Done", "duration", time.Since(t))
				}

				// if is watching, deploy process can run in goroutine
				if b.IsWatching() {
					go fn(name, gitOpt)
				} else {
					// otherwise, block it
					fn(name, gitOpt)
				}
			}

			return nil
		}
	}

	// build action
	return func(ctx *cli.Context) {
		if iniFile == nil {
			log15.Crit("Deploy.Init.Fail", "error", errors.New("please write deploy options to conf.ini"))
		}
		dOpt, err := deploy.NewOption(iniFile)
		if err != nil {
			log15.Error("Deploy.Init.Fail", "error", err)
			return
		}

		// add hook to opt
		if opt.After == nil {
			opt.After = afterFunc(dOpt)
		} else {
			opt.After = func(b *builder.Builder, ctx *builder.Context) error {
				if err := opt.After(b, ctx); err != nil {
					return err
				}
				return afterFunc(dOpt)(b, ctx)
			}
		}

		// run build site
		buildSite(opt, false)(ctx)
	}
}
