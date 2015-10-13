package boot

import (
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/src/core"
	"github.com/go-xiaohei/pugo/src/service"
	"gopkg.in/inconshreveable/log15.v2"
	"strconv"
	"time"
)

var (
	installCommand = cli.Command{
		Name:   "install",
		Usage:  "install pugo blog service with default configurations and data",
		Action: Install,
		Flags: []cli.Flag{
			cli.StringFlag{"port", "9899", "http server port", ""},
			cli.StringFlag{"domain", "localhost", "http server public domain", ""},
			cli.StringFlag{"dsn", "data/tidb", "tidb connection string", ""},
			cli.StringFlag{"user", "admin", "admin user name", ""},
			cli.StringFlag{"email", "admin@example.com", "admin user email", ""},
			cli.StringFlag{"password", "123456789", "admin user password", ""},
		},
	}
)

// install command action
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
		opt2 := newBootstrapInstallOption(ctx)
		if err := service.Call(service.Bootstrap.Install, opt2); err != nil {
			log15.Crit("Install.fail", "error", err)
		}
		log15.Info("Install.finish")
		return
	}
	i, _ := strconv.ParseInt(core.Cfg.Install, 10, 64)
	log15.Warn("Install.HadInstalled", "version", core.Cfg.Version, "installed", time.Unix(i, 0))
}

func newBootstrapInstallOption(ctx *cli.Context) service.BootstrapInstallOption {
	opt := service.BootstrapInstallOption{
		Port:          ctx.String("port"),
		Domain:        ctx.String("domain"),
		DbDSN:         ctx.String("dsn"),
		AdminUser:     ctx.String("user"),
		AdminEmail:    ctx.String("email"),
		AdminPassword: ctx.String("password"),
	}
	if opt.Port == "" {
		opt.Port = "9899"
	}
	if opt.Domain == "" {
		opt.Domain = "localhost"
	}
	if opt.DbDSN == "" {
		opt.DbDSN = "data/tidb"
	}
	if opt.AdminUser == "" {
		opt.AdminUser = "admin"
	}
	if opt.AdminEmail == "" {
		opt.AdminEmail = "admin@example.com"
	}
	if opt.AdminPassword == "" {
		opt.AdminPassword = "123456789"
	}
	return opt
}
