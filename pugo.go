package main

import (
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo-static/app/builder"
	"github.com/go-xiaohei/pugo-static/app/command"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/inconshreveable/log15.v2/ext"
	"path"
)

//go:generate go-bindata -o=app/asset/asset.go -pkg=asset source/... template/...

const (
	VERSION  = "0.8.0"
	VER_DATE = "2015-11-30"

	SRC_DIR   = "source"   // source contents dir
	TPL_DIR   = "template" // template dir
	MEDIA_DIR = "media"    // upload dir
	DST_DIR   = "dest"     // destination dir
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
	opt.SrcDir = SRC_DIR
	opt.TplDir = TPL_DIR
	opt.MediaDir = path.Join(SRC_DIR, MEDIA_DIR)
	opt.Version = VERSION
	opt.VerDate = VER_DATE

	log15.Root().SetHandler(log15.LvlFilterHandler(log15.LvlInfo, ext.FatalHandler(log15.StderrHandler)))
}

func main() {
	// app.Action = action
	app.Commands = []cli.Command{
		command.New(SRC_DIR, TPL_DIR),
		command.Build(opt),
		command.Server(opt),
		command.Deploy(opt),
	}
	app.RunAndExitOnError()
}
