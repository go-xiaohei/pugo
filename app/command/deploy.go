package command

import (
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo-static/app/builder"
	"github.com/go-xiaohei/pugo-static/app/deploy"
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
		},
		Action: deploySite(opt),
		Before: setDebugMode,
	}
}

func deploySite(opt *builder.BuildOption) func(ctx *cli.Context) {
	// build action
	return func(ctx *cli.Context) {
		deployer := deploy.New()
		// add hook to opt
		opt.After(func(b *builder.Builder, c *builder.Context) error {
            return deployer.Run(b,c)
        })
		// run build site
		buildSite(opt, false)(ctx)
	}
}
