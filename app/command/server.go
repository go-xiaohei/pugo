package command

import (
	"net/http"
	// pprof to profile
	_ "net/http/pprof"
	"time"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/builder"
	"github.com/go-xiaohei/pugo/app/server"
	"github.com/urfave/cli"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	// Server is command of 'server'
	Server = cli.Command{
		Name:  "server",
		Usage: "server static files",
		Flags: []cli.Flag{
			buildSourceFlag,
			buildDestFlag,
			buildThemeFlag,
			addrFlag,
			serveStaticFlag,
			debugFlag,
			noWatchFlag,
			cli.BoolFlag{
				Name: "profile",
			},
		},
		Before: Before,
		Action: serv,
	}

	s *server.Server
)

func serv(c *cli.Context) error {
	if c.Bool("static") {
		ctx := newContext(c, false)
		builder.Read(ctx)

		dstDir := ctx.DstDir()
		if !com.IsDir(dstDir) {
			log15.Crit("Server|Dest|'%s' is not directory", dstDir)
		}
		log15.Info("Server|Static|%s", dstDir)
		s := server.New(dstDir)
		s.SetPrefix(ctx.Source.Meta.Path)
		s.Run(c.String("addr"))
		return nil
	}
	builder.After(func(ctx *builder.Context) {
		if s == nil {
			s = server.New(ctx.DstDir())
			go s.Run(c.String("addr"))
		}
		if ctx.Source != nil && ctx.Source.Meta != nil {
			s.SetPrefix(ctx.Source.Meta.Path)
		}
	})

	if c.Bool("profile") {
		go http.ListenAndServe("localhost:6060", nil)
		for {
			build(newContext(c, false), false)
			time.Sleep(time.Second)
		}
	} else {
		if c.Bool("no-watch") {
			buildHangUp(newContext(c, true))
			return nil
		}
		build(newContext(c, true), true)
	}

	return nil
}
