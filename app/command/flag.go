package command

import (
	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/inconshreveable/log15.v2/ext"
	"gopkg.in/ini.v1"
)

var (
	// conf.ini
	confFile = "conf.ini"
	iniFile  *ini.File

	debugFlag = cli.BoolFlag{
		Name:  "debug",
		Usage: "print debug logs",
	}
	themeFlag = cli.StringFlag{
		Name:  "theme",
		Value: "default",
		Usage: "set theme to render",
	}
	destFlag = cli.StringFlag{
		Name:  "dest",
		Value: "dest",
		Usage: "set compiling to directory",
	}
	addrFlag = cli.StringFlag{
		Name:  "addr",
		Value: "0.0.0.0:9899",
		Usage: "set http server address",
	}
	watchFlag = cli.BoolFlag{
		Name:  "watch",
		Usage: "watch changes and auto-rebuild",
	}

	// global conf from conf.ini
	isDebug bool = false
	isWatch bool = false
)

func init() {
	if !com.IsFile(confFile) {
		return
	}
	iFile, err := ini.Load(confFile)
	if err != nil {
		log15.Crit("Conf.Load."+confFile, "error", err.Error())
		return
	}
	iniFile = iFile
	isDebug = iniFile.Section("mode").Key("debug").MustBool(false)
	isWatch = iniFile.Section("build").Key("watch").MustBool(false)

	themeFlag.Value = iniFile.Section("build").Key("theme").MustString("default")
	destFlag.Value = iniFile.Section("build").Key("dest").MustString("dest")
	addrFlag.Value = iniFile.Section("server").Key("addr").MustString("0.0.0.0:9899")
}

func setDebugMode(ctx *cli.Context) error {
	if ctx.Bool("debug") || isDebug {
		log15.Root().SetHandler(log15.LvlFilterHandler(log15.LvlDebug, ext.FatalHandler(log15.StderrHandler)))
		if com.IsFile(confFile) {
			log15.Debug("Conf.Load." + confFile)
		}
	}
	return nil
}
