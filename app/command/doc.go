package command

import (
	"github.com/go-xiaohei/pugo/app/builder"
	"github.com/go-xiaohei/pugo/app/server"
	"github.com/urfave/cli"
)

var (
	// Doc is command of 'doc'
	Doc = cli.Command{
		Name:  "doc",
		Usage: "run documentation server",
		Flags: []cli.Flag{
			addrFlag,
			debugFlag,
		},
		Before: Before,
		Action: docServ,
	}
)

func docServ(c *cli.Context) error {
	builder.After(func(ctx *builder.Context) {
		if s == nil {
			s = server.New(ctx.DstDir())
			go s.Run(c.String("addr"))
		}
		if ctx.Source != nil && ctx.Source.Meta != nil {
			s.SetPrefix(ctx.Source.Meta.Path)
		}
	})
	buildContext := newContext(c, false)
	buildContext.From = "doc/source"
	buildContext.To = "doc/dest"
	buildContext.ThemeName = "doc/theme"
	build(buildContext, true)
	return nil
}
