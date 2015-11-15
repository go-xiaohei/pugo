package main

import (
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo-static/builder"
	"github.com/go-xiaohei/pugo-static/command"
	"github.com/go-xiaohei/pugo-static/server"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/inconshreveable/log15.v2/ext"
)

//go:generate go-bindata -o=asset/asset.go -pkg=asset source/... template/...

const (
	VERSION  = "0.7.0"
	VER_DATE = "2015-11-14"

	SRC_DIR    = "source"   // source contents dir
	TPL_DIR    = "template" // template dir
	UPLOAD_DIR = "upload"   // upload dir
	DST_DIR    = "dest"     // destination dir
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
		cli.BoolFlag{
			Name:  "build",
			Usage: "only build site, but don't serve http",
		},
	}
	log15.Root().SetHandler(log15.LvlFilterHandler(log15.LvlDebug, ext.FatalHandler(log15.StderrHandler)))
}

func action(ctx *cli.Context) {
	if !ctx.Bool("debug") {
		log15.Root().SetHandler(log15.LvlFilterHandler(log15.LvlInfo, ext.FatalHandler(log15.StderrHandler)))
	}

	log15.Debug("Dir.Source./" + SRC_DIR)
	log15.Debug("Dir.Template./" + TPL_DIR)
	log15.Debug("Dir.Destination./" + DST_DIR)

	// builder
	b := builder.New(SRC_DIR, TPL_DIR, ctx.String("theme"), ctx.Bool("debug"))
	if b.Error != nil {
		panic(b.Error)
	}

	b.Build(DST_DIR)
	if ctx.Bool("debug") {
		b.Watch(DST_DIR, TPL_DIR)
	} else {
		b.Watch(DST_DIR, "")
	}

	if ctx.Bool("build") {
		return
	}

	// server
	staticDir := b.Renders().Current().StaticDir()
	static := server.NewStatic()
	static.RootPath = staticDir

	upload := server.NewStatic()
	upload.RootPath = UPLOAD_DIR
	upload.Prefix = "/upload"
	upload.IndexFiles = []string{}

	s := server.NewServer(ctx.String("addr"))
	s.Static = []*server.Static{static, upload}
	s.Helper = server.NewHelper(b, DST_DIR)
	s.ErrorHandler = server.Errors(DST_DIR)
	s.Run()
}

func main() {
	app.Action = action
	app.Commands = []cli.Command{command.New(SRC_DIR, TPL_DIR)}
	app.RunAndExitOnError()
}
