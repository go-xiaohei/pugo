package command

import (
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/builder"
	"github.com/go-xiaohei/pugo/app/server"
)

var (
	// Server is command of 'server'
	Server = cli.Command{
		Name:  "server",
		Usage: "server static files",
		Flags: []cli.Flag{
			buildFromFlag,
			buildToFlag,
			themeFlag,
			addrFlag,
			debugFlag,
		},
		Before: Before,
		Action: serv,
	}

	s *server.Server
)

func serv(c *cli.Context) {
	builder.After(func(ctx *builder.Context) {
		if s == nil {
			s = server.New(ctx.DstDir())
			go s.Run(c.String("addr"))
		}
		if ctx.Source != nil && ctx.Source.Meta != nil {
			s.SetPrefix(ctx.Source.Meta.Path)
		}
	})
	build(newContext(c, true), true)
}
