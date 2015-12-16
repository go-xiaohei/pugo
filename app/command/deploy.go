package command

import (
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
			cli.BoolFlag{
				Name:  "git",
				Usage: "deploy via git",
			},
			cli.StringFlag{
				Name:  "ftp",
				Usage: "deploy via ftp",
			},
			cli.StringFlag{
				Name:  "sftp",
				Usage: "deploy via sftp",
			},
		},
		Action: deploySite(opt),
		Before: setDebugMode,
	}
}

func deploySite(opt *builder.BuildOption) func(ctx *cli.Context) {
	// build action
	return func(ctx *cli.Context) {
		deployer := &deploy.Deployer{}

		// add git task
		tasks := []string{}
		if ctx.Bool("git") {
			tasks = append(tasks, "git://")
		}
		if ftpStr := ctx.String("ftp"); ftpStr != "" {
			tasks = append(tasks, "ftp:"+ftpStr)
		}
		if sftpStr := ctx.String("sftp"); sftpStr != "" {
			tasks = append(tasks, "sftp:"+sftpStr)
		}

		var err error
		for _, task := range tasks {
			if err = deployer.Add(task); err != nil {
				log15.Crit("Deploy.Fail", "error", err.Error())
			}
		}

		// add hook to opt
		opt.After(func(b *builder.Builder, c *builder.Context) error {
			return deployer.Run(b, c)
		})
		// run build site
		buildSite(opt, false)(ctx)
	}
}
