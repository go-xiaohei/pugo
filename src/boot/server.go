package boot

import (
	"github.com/codegangsta/cli"
	"gopkg.in/inconshreveable/log15.v2"
	"pugo/src/core"
	"pugo/src/service"
)

var (
	serverCommand = cli.Command{
		Name:   "server",
		Usage:  "start pugo blog web server",
		Action: Server,
	}
)

func Server(ctx *cli.Context) {
	opt := service.BootstrapOption{true, false, false}
	if err := service.Call(service.Bootstrap.Init, opt); err != nil {
		log15.Crit("Server.start.fail", "error", err)
	}
	// if not installed,try to install
	if core.Cfg.Install == "0" {
		Install(ctx)
	}

	// re-init all things
	opt = service.BootstrapOption{true, true, true}
	if err := service.Call(service.Bootstrap.Init, opt); err != nil {
		log15.Crit("Server.start.fail", "error", err)
	}
	log15.Info("Server.prepare")

	opt2 := service.UserAuthOption{"admin", "", "123456789", 3600, "webpage"}
	service.Call(service.User.Authorize, opt2)

	// start server
	log15.Info("Server.start." + core.Cfg.Http.Host + ":" + core.Cfg.Http.Port)
	core.Server.Run(core.Cfg.Http.Host + ":" + core.Cfg.Http.Port)
}
