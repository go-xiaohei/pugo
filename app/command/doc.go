package command

import (
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/builder"
)

// Server command serve files
func Doc(opt *builder.BuildOption) cli.Command {
	return cli.Command{
		Name:     "doc",
		Usage:    "build and serve documentation",
		HideHelp: true,
		Flags: []cli.Flag{
			addrFlag,
			debugFlag,
		},
		Action: func() func(*cli.Context) {
			// clone option
			opt2 := *opt
			opt2.SrcDir = "doc/source"
			opt2.TplDir = "doc/theme"
			opt2.Theme = "default"
			opt2.TargetDir = "doc/dest"
			return serveSite(&opt2)
		}(),
		Before: setDebugMode,
	}
}
