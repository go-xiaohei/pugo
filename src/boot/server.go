package boot

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/pugo/src/controller/admin"
	"github.com/fuxiaohei/pugo/src/controller/public"
	"github.com/fuxiaohei/pugo/src/core"
	"github.com/fuxiaohei/pugo/src/middle"
	"github.com/fuxiaohei/pugo/src/service"
	"github.com/lunny/tango"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	serverCommand = cli.Command{
		Name:   "server",
		Usage:  "start pugo blog web server",
		Action: Server,
	}
)

func Server(ctx *cli.Context) {
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
	core.Server.Get("/feed.xml", new(public.RssController))
	core.Server.Get("/archive", new(public.ArchiveController))
	core.Server.Get("/:link.html", new(public.PageController))
	core.Server.Get("/", new(public.IndexController))

	// start server
	log15.Info("Server.start." + core.Cfg.Http.Host + ":" + core.Cfg.Http.Port)
	core.Server.Run(core.Cfg.Http.Host + ":" + core.Cfg.Http.Port)
}
