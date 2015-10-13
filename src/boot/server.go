package boot

import (
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/src/controller/admin"
	"github.com/go-xiaohei/pugo/src/controller/public"
	"github.com/go-xiaohei/pugo/src/core"
	"github.com/go-xiaohei/pugo/src/middle"
	"github.com/go-xiaohei/pugo/src/service"
	"github.com/lunny/tango"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/inconshreveable/log15.v2/ext"
)

var (
	serverCommand = cli.Command{
		Name:  "server",
		Usage: "start pugo blog web server",
		Action: func(ctx *cli.Context) {
			sc := &serverContext{ctx: ctx}
			core.Start(sc)
		},
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "debug",
				Usage: "show more debug info when running server",
			},
		},
	}
)

type serverContext struct {
	ctx *cli.Context
}

func (sc *serverContext) Start() {
	Server(sc.ctx)
}

func (sc *serverContext) Stop() {

}

// server command action,
// it's used by *boot.serverContext when run "pugo.exe server" command
func Server(ctx *cli.Context) {
	// change logger level
	if ctx.Bool("debug") {
		core.RunMode = core.RUN_MOD_DEBUG
		log15.Root().SetHandler(ext.FatalHandler(log15.StderrHandler))
	}

	opt := service.BootstrapInitOption{true, false, false}
	if err := service.Call(service.Bootstrap.Init, opt); err != nil {
		log15.Crit("Server.start.fail", "error", err)
	}
	// if not installed,try to install
	if core.Cfg.Install == "0" {
		Install(ctx)
	}

	// re-init all things
	opt = service.BootstrapInitOption{true, true, true}
	if err := service.Call(service.Bootstrap.Init, opt); err != nil {
		log15.Crit("Server.start.fail", "error", err)
	}
	// bootstrap service, preload data
	opt2 := service.BootstrapOption{true, true, true}
	if err := service.Call(service.Bootstrap.Bootstrap, opt2); err != nil {
		log15.Crit("Server.start.fail", "error", err)
	}
	log15.Info("Server.prepare")

	// set middleware and routers
	core.Server.ErrHandler = new(middle.RecoveryHandler)
	core.Server.Use(
		middle.Recover(),
		middle.Logger(),
		middle.Themer(),
		middle.Responser(),
		middle.Authorizor())

	adminGroup := tango.NewGroup()
	adminGroup.Any("/login", new(admin.LoginController))
	adminGroup.Route([]string{"GET:Logout"}, "/logout", new(admin.LoginController))

	adminGroup.Any("/write/article", new(admin.ArticleWriteController))
	adminGroup.Get("/manage/article", new(admin.ArticleManageController))
	adminGroup.Get("/public/article", new(admin.ArticlePublicController))
	adminGroup.Get("/delete/article", new(admin.ArticleDeleteController))

	adminGroup.Any("/write/page", new(admin.PageWriteController))
	adminGroup.Get("/manage/page", new(admin.PageManageController))
	adminGroup.Get("/delete/page", new(admin.PageDeleteController))

	adminGroup.Get("/manage/comment", new(admin.CommentController))
	adminGroup.Route([]string{"GET:Approve"}, "/approve/comment", new(admin.CommentController))
	adminGroup.Route([]string{"GET:Delete"}, "/delete/comment", new(admin.CommentController))
	adminGroup.Route([]string{"POST:Reply"}, "/reply/comment", new(admin.CommentController))

	adminGroup.Any("/profile", new(admin.ProfileController))
	adminGroup.Route([]string{"POST:Password"}, "/password", new(admin.ProfileController))

	adminGroup.Any("/option/general", new(admin.SettingGeneralController))
	adminGroup.Route([]string{"POST:PostMedia"}, "/option/media", new(admin.SettingGeneralController))
	adminGroup.Get("/option/theme", new(admin.SettingThemeController))
	adminGroup.Any("/option/content", new(admin.SettingContentController))
	adminGroup.Any("/option/comment", new(admin.SettingCommentController))
	adminGroup.Any("/option/menu", new(admin.SettingMenuController))

	adminGroup.Get("/manage/media", new(admin.MediaController))
	adminGroup.Route([]string{"POST:Upload"}, "/upload/media", new(admin.MediaController))
	adminGroup.Get("/delete/media", new(admin.MediaDeleteController))

	adminGroup.Get("/advance/backup", new(admin.AdvBackupController))
	adminGroup.Route([]string{"POST:Backup"}, "/advance/backup", new(admin.AdvBackupController))
	adminGroup.Route([]string{"GET:Delete"}, "/delete/backup", new(admin.AdvBackupController))

	adminGroup.Get("/advance/import", new(admin.AdvImportController))
	adminGroup.Post("/import/:type", new(admin.AdvImportController))

	adminGroup.Get("/", new(admin.IndexController))
	core.Server.Group("/admin", adminGroup)

	core.Server.Get("/article/page/:page", new(public.IndexController))
	core.Server.Get("/article/:id/:link.html", new(public.ArticleController))
	core.Server.Get("/page/:id/:link.html", new(public.PageController))
	core.Server.Post("/comment/:type/:id", new(public.CommentController))

	core.Server.Get("/robots.txt", new(public.RobotController))
	core.Server.Get("/feed.xml", new(public.RssController))
	core.Server.Get("/feed", new(public.RssController))

	core.Server.Get("/archive", new(public.ArchiveController))
	core.Server.Get("/:link.html", new(public.PageController))

	core.Server.Route([]string{"GET:NotFound"}, "/error/404.html", new(public.ErrorController))
	core.Server.Route([]string{"GET:InternalError"}, "/error/503.html", new(public.ErrorController))

	core.Server.Get("/", new(public.IndexController))

	// start server
	addr := core.Cfg.Http.Host + ":" + core.Cfg.Http.Port
	log15.Info("Server.start." + addr)
	go core.Server.Run(addr)
}
