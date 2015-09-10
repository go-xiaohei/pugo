package boot

import (
	"github.com/codegangsta/cli"
	"gopkg.in/inconshreveable/log15.v2"
	"pugo/src/core"
	"pugo/src/service"
	"strconv"
	"time"
)

var (
	installCommand = cli.Command{
		Name:   "install",
		Usage:  "install pugo blog service with default configurations and data",
		Action: Install,
	}
)

func Install(ctx *cli.Context) {
	opt := service.BootstrapInitOption{true, false, false}
	if err := service.Call(service.Bootstrap.Init, opt); err != nil {
		log15.Crit("Install.fail", "error", err)
	}
	if core.Cfg.Install == "0" {
		log15.Info("Install.start")
		opt = service.BootstrapInitOption{false, true, false} // connect to database
		if err := service.Call(service.Bootstrap.Init, opt); err != nil {
			log15.Crit("Install.fail", "error", err)
		}
		if err := service.Call(service.Bootstrap.Install, nil); err != nil {
			log15.Crit("Install.fail", "error", err)
		}
		log15.Info("Install.finish")
		return
	}
	i, _ := strconv.ParseInt(core.Cfg.Install, 10, 64)
	log15.Warn("Install.HadInstalled", "version", core.Cfg.Version, "installed", time.Unix(i, 0))
}
