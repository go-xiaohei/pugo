package command

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo-static/app/builder"
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

type deployOption struct {
}

func deploySite(opt *builder.BuildOption) func(ctx *cli.Context) {
	// real deploy action, in builder hook
	afterFunc := func(dopt *deployOption) builder.BuildHook {
		return func(b *builder.Builder, ctx *builder.Context) error {
			println("deploy")
			return nil
		}
	}

	// build action
	return func(ctx *cli.Context) {
		if iniFile == nil {
			log15.Crit("Deploy.Fail", "error", errors.New("please write deploy options to conf.ini"))
		}
		fmt.Println(confFile)

		dOpt := &deployOption{}

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
