package main

import (
	"github.com/codegangsta/cli"
	"pugo/builder"
	"pugo/server"
)

const (
	VERSION  = "1.0"
	VER_DATE = "2015-11-05"
	SRC_DIR  = "source"
	TPL_DIR  = "template"
)

var (
	app = cli.NewApp()
)

func init() {
	app.Name = "pugo"
	app.Usage = "a beautiful site generator"
	app.Author = "https://github.com/fuxiaohei"
	app.Email = "fuxiaohei@vip.qq.com"
	app.Version = VERSION + "(" + VER_DATE + ")"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Value: "0.0.0.0:9899",
			Usage: "pugo's http server address",
		},
		cli.StringFlag{
			Name:  "theme",
			Value: "default",
			Usage: "pugo's theme to display",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "debug mode flag",
		},
	}
	app.Action = action
}

func action(ctx *cli.Context) {
	// builder
	b := builder.New(SRC_DIR, TPL_DIR, ctx.String("theme"), ctx.Bool("debug"))
	if b.Error != nil {
		panic(b.Error)
	}

	go b.Build()

	// server
	staticDir := b.Renders().Current().StaticDir()
	static := server.NewStatic()
	static.RootPath = staticDir
	s := server.NewServer(ctx.String("addr"), static, server.NewHelper(b))
	s.Run()
}

func main() {
	app.RunAndExitOnError()
}
