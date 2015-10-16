package service

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xiaohei/pugo/src/core"
	"github.com/go-xiaohei/pugo/src/model"
	"github.com/go-xiaohei/pugo/src/utils"
	_ "github.com/go-xorm/tidb"
	"github.com/go-xorm/xorm"
	"github.com/lunny/tango"
	"github.com/ngaut/log"
	"github.com/pingcap/tidb"
	"github.com/tango-contrib/binding"
	"github.com/tango-contrib/flash"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/session"
	"github.com/tango-contrib/xsrf"
	"html/template"
	"net/http"
	_ "net/http/pprof"
	"path"
	"time"
)

var (
	Bootstrap = new(BootstrapService)
)

type BootstrapInitOption struct {
	Config   bool
	Database bool
	Server   bool
}

type BootstrapService struct{}

func (is *BootstrapService) Init(v interface{}) (*Result, error) {
	opt, ok := v.(BootstrapInitOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(is.Init, opt)
	}
	var err error
	if opt.Config {
		core.Cfg = core.NewConfig()
		if err = core.Cfg.Sync(core.ConfigFile); err != nil {
			return nil, err
		}
	}
	if core.Cfg != nil && opt.Database { // database depends on config
		log.SetLevelByString("error")
		tidb.Debug = false // close tidb debug
		core.Db, err = xorm.NewEngine(core.Cfg.Db.Driver, core.Cfg.Db.DSN)
		if err != nil {
			return nil, err
		}
		// debug mode
		if core.RunMode == core.RUN_MOD_DEBUG {
			core.Db.ShowDebug = true
			core.Db.ShowInfo = false
			core.Db.ShowSQL = true
			core.Db.ShowWarn = true
			core.Db.ShowErr = true
			go func() {
				http.ListenAndServe("0.0.0.0:6060", nil)
			}()
		} else {
			// close log in prod mode
			core.Db.SetLogger(nil)
			cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
			core.Db.SetDefaultCacher(cacher)
		}

		// ping to test connection
		if err := core.Db.Ping(); err != nil {
			return nil, err
		}
	}
	if core.Cfg != nil && opt.Server { // server depends on config
		core.Server = tango.New([]tango.Handler{
			tango.Return(),
			tango.Param(),
			tango.Contexts(),
		}...)
		core.Server.Use(tango.Static(tango.StaticOptions{
			RootPath: core.StaticDirectory,
			Prefix:   core.StaticPrefix,
		}))
		core.Server.Use(tango.Static(tango.StaticOptions{
			RootPath: core.ThemeDirectory,
			Prefix:   core.ThemePrefix,
		}))
		core.Render = renders.New(renders.Options{
			Reload:     true,
			Directory:  core.ThemeDirectory,
			Extensions: []string{".tmpl"},
			Funcs: template.FuncMap{
				"TimeUnixFormat":  utils.TimeUnixFormat,
				"TimeUnixFriend":  utils.TimeUnixFriend,
				"Mardown2Str":     utils.Markdown2String,
				"Markdown2HTML":   utils.Markdown2HTML,
				"Nl2BrHTML":       utils.Nl2Br,
				"Nl2BrString":     utils.Nl2BrString,
				"BytesSizeFriend": utils.FriendBytesSize,
			},
		})
		core.Server.Use(core.Render)
		core.Server.Use(xsrf.New(time.Hour))
		core.Server.Use(binding.Bind())
		core.Session = session.New(session.Options{
			SessionIdName: core.SessionName,
		})
		core.Server.Use(core.Session)
		core.Server.Use(flash.Flashes(core.Session))
	}
	return nil, nil
}

type BootstrapInstallOption struct {
	Port          string
	Domain        string
	DbDSN         string
	AdminUser     string
	AdminEmail    string
	AdminPassword string
}

func (bs *BootstrapService) Install(v interface{}) (*Result, error) {
	opt, ok := v.(BootstrapInstallOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(bs.Install, opt)
	}

	// create tables
	if err := core.Db.Sync2(new(model.User),
		new(model.UserToken),
		new(model.Theme),
		new(model.Article),
		new(model.ArticleTag),
		new(model.Setting),
		new(model.Media),
		new(model.Page),
		new(model.Comment),
		new(model.Message)); err != nil {
		return nil, err
	}

	// insert default user
	user := &model.User{
		Name:      opt.AdminUser,
		Email:     opt.AdminEmail,
		Nick:      opt.AdminUser,
		Profile:   "this is administrator",
		Role:      model.USER_ROLE_ADMIN,
		Status:    model.USER_STATUS_ACTIVE,
		AvatarUrl: utils.Gravatar("admin@admin.com"),
	}
	user.SetPassword(opt.AdminPassword)
	if _, err := core.Db.Insert(user); err != nil {
		return nil, err
	}

	// insert default themes
	themes := []interface{}{
		&model.Theme{
			Name:      "admin",
			Author:    core.PUGO_AUTHOR,
			Version:   "1.0",
			Directory: path.Join(core.ThemeDirectory, "admin"),
			Status:    model.THEME_STATUS_LOCKED,
		},
		&model.Theme{
			Name:      "default",
			Author:    core.PUGO_AUTHOR,
			Version:   "1.0",
			Directory: path.Join(core.ThemeDirectory, "default"),
			Status:    model.THEME_STATUS_CURRENT,
		},
	}
	if _, err := core.Db.Insert(themes...); err != nil {
		return nil, err
	}

	// insert settings
	generalSetting := &model.SettingGeneral{
		Title:          "PUGO",
		SubTitle:       "Simple Blog Engine",
		Keyword:        "pugo,blog,go,golang",
		Description:    "PUGO is a simple blog engine by golang",
		HostName:       fmt.Sprintf("http://%s", opt.Domain),
		HeroImage:      "/img/bg.png",
		TopAvatarImage: "/img/logo.png",
	}
	setting := &model.Setting{
		Name:   "general",
		UserId: 0,
		Type:   model.SETTING_TYPE_GENERAL,
	}
	setting.Encode(generalSetting)
	if _, err := core.Db.Insert(setting); err != nil {
		return nil, err
	}
	Setting.General = generalSetting

	mediaSetting := &model.SettingMedia{
		MaxFileSize: 10 * 1024,
		ImageFile:   []string{"jpg", "jpeg", "png", "gif", "bmp", "vbmp"},
		DocFile:     []string{"txt", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "pdf"},
		CommonFile:  []string{"zip", "rar"},
		DynamicLink: false,
	}
	setting = &model.Setting{
		Name:   "media",
		UserId: 0,
		Type:   model.SETTING_TYPE_MEDIA,
	}
	setting.Encode(mediaSetting)
	if _, err := core.Db.Insert(setting); err != nil {
		return nil, err
	}
	Setting.Media = mediaSetting

	contentSetting := &model.SettingContent{
		PageSize:         5,
		RSSFullText:      true,
		RSSNumberLimit:   0,
		TopPage:          0,
		PageDisallowLink: []string{"article", "archive", "feed", "comment", "admin", "sitemap"},
	}
	setting = &model.Setting{
		Name:   "content",
		UserId: 0,
		Type:   model.SETTING_TYPE_CONTENT,
	}
	setting.Encode(contentSetting)
	if _, err := core.Db.Insert(setting); err != nil {
		return nil, err
	}
	Setting.Content = contentSetting

	commentSetting := &model.SettingComment{
		IsPager:        false,
		PageSize:       10,
		Order:          "create_time DESC",
		CheckAll:       false,
		CheckNoPass:    true,
		CheckRefer:     true,
		AutoCloseDay:   30,
		SubmitDuration: 60,
		MaxLength:      512,
		MinLength:      2,
	}
	setting = &model.Setting{
		Name:   "comment",
		UserId: 0,
		Type:   model.SETTING_TYPE_COMMENT,
	}
	setting.Encode(commentSetting)
	if _, err := core.Db.Insert(setting); err != nil {
		return nil, err
	}
	Setting.Comment = commentSetting

	menuSettings := []*model.SettingMenu{
		{
			"Home", "/", "Home",
			false, "home",
		},
		{
			"Archive", "/archive", "Archive",
			false, "archive",
		},
		{
			"About", "/about.html", "About",
			false, "about",
		},
	}
	setting = &model.Setting{
		Name:   "menu",
		UserId: 0,
		Type:   model.SETTING_TYPE_MENU,
	}
	setting.Encode(menuSettings)
	if _, err := core.Db.Insert(setting); err != nil {
		return nil, err
	}
	Setting.Menu = menuSettings

	// first article
	article := &model.Article{
		UserId:        user.Id,
		Title:         firstArticleTitle,
		Link:          firstArticleLink,
		Body:          firstArticleContent,
		TagString:     firstArticleTag,
		Status:        model.ARTICLE_STATUS_PUBLISH,
		CommentStatus: model.ARTICLE_COMMENT_OPEN,
		Hits:          1,
		Preview:       firstArticleContent,
		BodyType:      model.ARTICLE_BODY_MARKDOWN,
	}
	if _, err := Article.Write(article); err != nil {
		return nil, err
	}

	// first comment
	cmt := &model.Comment{
		Name:      user.Name,
		UserId:    user.Id,
		Email:     user.Email,
		Url:       user.Url,
		AvatarUrl: user.AvatarUrl,
		Body:      firstCommentContent,
		Status:    model.COMMENT_STATUS_APPROVED,
		From:      model.COMMENT_FROM_ARTICLE,
		FromId:    article.Id,
		ParentId:  0,
	}
	if _, err := Comment.Save(cmt); err != nil {
		return nil, err
	}

	// first page
	page := &model.Page{
		UserId:        user.Id,
		Title:         firstPageTitle,
		Link:          firstPageLink,
		Body:          firstPageContent,
		Status:        model.PAGE_STATUS_PUBLISH,
		CommentStatus: model.PAGE_COMMENT_OPEN,
		Hits:          1,
		Template:      "page.tmpl",
		BodyType:      model.PAGE_BODY_MARKDOWN,
		TopLink:       true,
	}
	if _, err := Page.Write(page); err != nil {
		return nil, err
	}

	// assign install time to config
	core.Cfg.Http.Port = opt.Port
	core.Cfg.Http.Domain = opt.Domain
	core.Cfg.Install = fmt.Sprint(time.Now().Unix())
	if err := core.Cfg.WriteToFile(core.ConfigFile); err != nil {
		return nil, err
	}
	return nil, nil
}

type BootstrapOption struct {
	Themes  bool // load themes
	I18n    bool // load languages
	Setting bool // load settings to SettingService
}

// bootstrap means loading memory data and starting some worker in background
func (bs *BootstrapService) Bootstrap(v interface{}) (*Result, error) {
	opt, ok := v.(BootstrapOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(bs.Bootstrap, opt)
	}
	if opt.Themes {
		if err := Call(Theme.Load, nil); err != nil {
			return nil, err
		}
	}
	if opt.Setting {
		var (
			sOpt    = SettingReadOption{model.SETTING_TYPE_GENERAL, 0, false}
			setting = new(model.Setting)
		)
		if err := Call(Setting.Read, sOpt, setting); err != nil {
			return nil, err
		}
		Setting.General = setting.ToGeneral()

		sOpt = SettingReadOption{model.SETTING_TYPE_MEDIA, 0, false}
		setting = new(model.Setting)
		if err := Call(Setting.Read, sOpt, setting); err != nil {
			return nil, err
		}
		Setting.Media = setting.ToMedia()

		sOpt = SettingReadOption{model.SETTING_TYPE_CONTENT, 0, false}
		setting = new(model.Setting)
		if err := Call(Setting.Read, sOpt, setting); err != nil {
			return nil, err
		}
		Setting.Content = setting.ToContent()

		sOpt = SettingReadOption{model.SETTING_TYPE_COMMENT, 0, false}
		setting = new(model.Setting)
		if err := Call(Setting.Read, sOpt, setting); err != nil {
			return nil, err
		}
		Setting.Comment = setting.ToComment()

		sOpt = SettingReadOption{model.SETTING_TYPE_MENU, 0, false}
		setting = new(model.Setting)
		if err := Call(Setting.Read, sOpt, setting); err != nil {
			return nil, err
		}
		Setting.Menu = setting.ToMenu()
	}
	return nil, nil
}

var (
	firstArticleTitle   = "Hello World"
	firstArticleLink    = "hello-world"
	firstArticleTag     = "hello"
	firstArticleContent = `# Hello World

Welcome to [Pugo](http://github.com/go-xiaohei/pugo)! This is your very first article. Read the [Wiki](http://github.com/go-xiaohei/pugo/wiki) for more infomation. If you get any problems when trying Pugo, you can find the answer or make a question in [issues](http://github.com/go-xiaohei/pugo/issues).

### Usage

You can sign in [admin panel](/admin/) with ` + "`admin`" + ` & ` + "`123456789`" + ` to change settings, write new article or page and upload media file.

**You'd better change default user and password setting to your own when first run.**

### Thanks

 - [xorm](https://github.com/go-xorm/xorm)
 - [tango](https://github.com/lunny/tango)
 - [tidb](https://github.com/pingcap/tidb)
 - [editor.md](https://github.com/pandao/editor.md)

`
	firstCommentContent = "this is first comment from administrator"

	firstPageTitle   = "About Pugo"
	firstPageLink    = "about"
	firstPageContent = "`Pugo`" + ` is a pure go blog engine to make new site. It works on ` + "`NewSQL`" + ` [tidb](https://github.com/pingcap/tidb) as an experiment. You write [Markdown]() content as an article or page with beautiful theme.


### Usage

You can download binary file from [Github Releases](https://github.com/go-xiaohei/pugo/releases) in your operation system.

Then unzip compressed file and run ` + "`pugo[.exe] server`" + ` to install and run site in ` + "`http://localhost:9899`" + `.

You need change ` + "`admin`" + ` settings to keep more safe.

### Contribute

Please feedback any question to [Github Issue](https://github.com/go-xiaohei/pugo/issues).
`
)
