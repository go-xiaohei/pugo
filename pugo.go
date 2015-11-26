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
	VERSION  = "0.7.5"
	VER_DATE = "2015-11-20"

	SRC_DIR    = "source"   // source contents dir
	TPL_DIR    = "template" // template dir
	UPLOAD_DIR = "upload"   // upload dir
	DST_DIR    = "dest"     // destination dir
)

var (
	app = cli.NewApp()
	opt = new(builder.BuildOption)
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
	opt.SrcDir = SRC_DIR
	opt.TplDir = TPL_DIR
	opt.UploadDir = UPLOAD_DIR
	opt.Version = VERSION
	opt.VerDate = VER_DATE

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
	opt.IsDebug = ctx.Bool("debug")
	if opt.IsDebug {
		opt.IsWatchTemplate = true
	}
	opt.Theme = ctx.String("theme")
	b := builder.New(opt)
	if b.Error != nil {
		panic(b.Error)
	}

	b.Build(DST_DIR)
	if ctx.Bool("build") {
		return
	}
	b.Watch(DST_DIR)

	// server
	static := server.NewStatic()
	static.RootPath = b.Context().Theme.Static() // use built context

	upload := server.NewStatic()
	upload.RootPath = UPLOAD_DIR
	upload.Prefix = "/upload"
	upload.IndexFiles = []string{}

	s := server.NewServer(ctx.String("addr"), b)
	s.Static = []*server.Static{static, upload}
	s.Helper = server.NewHelper(b, DST_DIR)
	s.ErrorHandler = server.Errors(DST_DIR)
	s.Run()
}

func main() {
	opt := &builder.BuildOption{
		SrcDir:    SRC_DIR,
		TplDir:    TPL_DIR,
		UploadDir: UPLOAD_DIR,
		Version:   VERSION,
		VerDate:   VER_DATE,
	}
	// app.Action = action
	app.Commands = []cli.Command{
		command.New(SRC_DIR, TPL_DIR),
		command.Build(opt),
		command.Server(),
	}
	app.RunAndExitOnError()
}
