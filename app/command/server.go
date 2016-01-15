package command

import (
	"net/http"

	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/builder"
	"github.com/go-xiaohei/pugo/app/server"
	"gopkg.in/inconshreveable/log15.v2"
)

// Server returns command serve files
func Server(opt *builder.Option) cli.Command {
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

func serveSite(opt *builder.Option) func(ctx *cli.Context) {
	return func(ctx *cli.Context) {
		if ctx.Bool("debug") {
			go http.ListenAndServe("0.0.0.0:6060", nil)
		}

		// set target dir
		if targetDir := ctx.String("dest"); targetDir != "" {
			opt.TargetDir = targetDir
		}
		if opt.TargetDir == "template" || opt.TargetDir == "source" {
			log15.Crit("Builder.Fail", "error", "destination directory should not be 'template' or 'source'")
		}

		s := server.New(opt.TargetDir)
		opt.After(func(b *builder.Builder, ctx *builder.Context) error {
			s.SetPrefix(ctx.Meta.Base)
			log15.Debug("Server.Prefix." + s.GetPrefix())
			return nil
		})

		// run server in goroutine
		go func() {
			addr := ctx.String("addr")
			s.Run(addr)
		}()
		// run buildSite to build
		buildSite(opt, true)(ctx)
	}
}
