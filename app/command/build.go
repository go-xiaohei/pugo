package command

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/builder"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	// Build is command of 'build'
	Build = cli.Command{
		Name:  "build",
		Usage: "build static files",
		Flags: []cli.Flag{
			buildFromFlag,
			buildToFlag,
			themeFlag,
			watchFlag,
			debugFlag,
		},
		Before: Before,
		Action: func(ctx *cli.Context) {
			build(ctx, false)
		},
	}
)

func build(c *cli.Context, mustWatch bool) {
	// ctrl+C capture
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	ctx := builder.NewContext(
		c.String("from"),
		c.String("to"),
		c.String("theme"),
	)
	if !ctx.IsValid() {
		log15.Crit("Build|Must have values in 'from', 'to' & 'theme'")
	}
	builder.Build(ctx)

	if c.Bool("watch") || mustWatch {
		builder.Watch(ctx)
		<-signalChan
		log15.Info("Watch|Close")
	}
}
