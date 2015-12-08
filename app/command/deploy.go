package command

import (
	"errors"

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

	if iniFile == nil {
		log15.Crit("Deploy.Fail", "error", errors.New("please add conf.ini to set deploy options"))
	}
	deployer, err := deploy.New(iniFile.Section("deploy"))
	if err != nil {
		log15.Crit("Deploy.Fail", "error", err.Error())
	}

	// real deploy action, in builder hook
	afterFunc := func(b *builder.Builder, ctx *builder.Context) error {
		if b.IsWatching() {
			return deployer.RunAsync()
		}
		return deployer.Run()
	}

	// build action
	return func(ctx *cli.Context) {

		// add hook to opt
		if opt.After == nil {
			opt.After = afterFunc
		} else {
			opt.After = func(b *builder.Builder, ctx *builder.Context) error {
				if err := opt.After(b, ctx); err != nil {
					return err
				}
				return afterFunc(b, ctx)
			}
		}

		// run build site
		buildSite(opt, false)(ctx)
	}
}
