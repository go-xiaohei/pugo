package main

import (
	_ "net/http/pprof"
	"path"

	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/builder"
	"github.com/go-xiaohei/pugo/app/command"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/inconshreveable/log15.v2/ext"
)

//go:generate go-bindata -o=app/asset/asset.go -pkg=asset source/... theme/... doc/source/... doc/theme/...

const (
	// Version number
	Version = "0.9.0.0103"

	// SrcDir contains contents
	SrcDir = "source"
	// ThemeDir contains themes
	ThemeDir = "theme"
	// MediaDir saves upload media dir
	MediaDir = "media"
)

var (
	app = cli.NewApp()
	opt = new(builder.Option)
)

func init() {
	app.Name = "PuGo"
	app.Usage = "a static website generator & deployer in Go"
	app.Author = "fuxiaohei"
	app.Email = "fuxiaohei@vip.qq.com"
	app.Version = Version
	opt.SrcDir = SrcDir
	opt.TplDir = ThemeDir
	opt.MediaDir = path.Join(SrcDir, MediaDir)
	opt.Version = Version

	log15.Root().SetHandler(log15.LvlFilterHandler(log15.LvlInfo, ext.FatalHandler(log15.StderrHandler)))
}

func main() {
	// app.Action = action
	app.Commands = []cli.Command{
		command.New(SrcDir, ThemeDir),
		command.Build(opt),
		command.Server(opt),
		command.Doc(opt),
	}
	app.RunAndExitOnError()
}
