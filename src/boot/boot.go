// Package boot provides bootstrap function
package boot

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/pugo/src/core"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/inconshreveable/log15.v2/ext"
	"os"
	"time"
)

func init() {
	// set log settings
	log15.Root().SetHandler(ext.FatalHandler(log15.StderrHandler))

	// set cli app
	core.App = cli.NewApp()
	core.App.Name = core.PUGO_NAME
	core.App.Usage = core.PUGO_DESCRIPTION
	core.App.Version = fmt.Sprintf("%s(%s)", core.PUGO_VERSION, core.PUGO_VERSION_STATE)
	core.App.Compiled, _ = time.Parse("20060102", core.PUGO_VERSION_DATE)
	core.App.Authors = []cli.Author{
		cli.Author{core.PUGO_AUTHOR, core.PUGO_AUTHOR_EMAIL},
	}
	core.App.HideHelp = true
	core.App.HideVersion = true
	core.App.CommandNotFound = func(_ *cli.Context, command string) {
		log15.Crit("command '" + command + "' is not found. please run 'pugo help'")
	}
	core.App.Commands = []cli.Command{
		installCommand,
		serverCommand,
	}

	// set crash log
	file, err := os.OpenFile(core.CrashLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log15.Crit("create crash log file error : " + err.Error())
	}
	core.Crash = log15.New()
	core.Crash.SetHandler(log15.StreamHandler(file, log15.JsonFormat()))
}

func Run() {
	core.App.RunAndExitOnError()
}
