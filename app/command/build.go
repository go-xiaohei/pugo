package command

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-xiaohei/pugo/app/builder"
	"github.com/urfave/cli"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	// Build is command of 'build'
	Build = cli.Command{
		Name:  "build",
		Usage: "build static files",
		Flags: []cli.Flag{
			buildSourceFlag,
			buildDestFlag,
			buildThemeFlag,
			buildWatchFlag,
			debugFlag,
		},
		Before: Before,
		Action: func(ctx *cli.Context) error {
			// migrate.Init()
			build(newContext(ctx, true), false)
			return nil
		},
	}
)

func newContext(c *cli.Context, validate bool) *builder.Context {
	ctx := builder.NewContext(
		c,
		c.String("source"),
		c.String("dest"),
		c.String("theme"),
	)
	if validate && !ctx.IsValid() {
		log15.Crit("Build|Must have values in 'source', 'dest' & 'theme'")
	}
	return ctx
}

func build(ctx *builder.Context, mustWatch bool) {
	// ctrl+C capture
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	builder.Build(ctx)

	if ctx.Cli().Bool("watch") || mustWatch {
		builder.Watch(ctx)
		<-signalChan
		log15.Info("Watch|Close")
	}
}
