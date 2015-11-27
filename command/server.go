package command

import (
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo-static/builder"
	"github.com/go-xiaohei/pugo-static/server"
)

func Server(opt *builder.BuildOption) cli.Command {
	return cli.Command{
		Name:     "server",
		Usage:    "build source and server static files",
		HideHelp: true,
		Flags: []cli.Flag{
			destFlag,
			themeFlag,
			addrFlag,
			debugFlag,
		},
		Action: serveSite(opt),
		Before: setDebugMode,
	}
}

func serveSite(opt *builder.BuildOption) func(ctx *cli.Context) {
	return func(ctx *cli.Context) {
		// run server in goroutine
		go serve(ctx, opt)
		// run buildSite to build
		buildSite(opt)(ctx)
	}
}

func serve(ctx *cli.Context, opt *builder.BuildOption) {
	s := server.New(ctx.String("dest"))
	addr := ctx.String("addr")
	s.Run(addr)
}
