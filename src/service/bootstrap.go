package service

import (
	"fmt"
	"github.com/fuxiaohei/pugo/src/core"
	"github.com/fuxiaohei/pugo/src/model"
	"github.com/fuxiaohei/pugo/src/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/lunny/tango"
	"github.com/tango-contrib/binding"
	"github.com/tango-contrib/flash"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/session"
	"github.com/tango-contrib/xsrf"
	"html/template"
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
		core.Db, err = xorm.NewEngine(core.Cfg.Db.Driver, core.Cfg.Db.DSN)
		if err != nil {
			return nil, err
		}
		core.Db.ShowDebug = true
		core.Db.ShowSQL = true
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
		core.Server.Use(renders.New(renders.Options{
			Reload:     true,
			Directory:  core.ThemeDirectory,
			Extensions: []string{".tmpl"},
			Funcs: template.FuncMap{
				"TimeUnixFormat": utils.TimeUnixFormat,
				"TimeUnixFriend": utils.TimeUnixFriend,
				"Mardown2Str":    utils.Markdown2String,
				"Markdown2HTML":  utils.Markdown2HTML,
				"Nl2BrHTML":      utils.Nl2Br,
				"Nl2BrString":    utils.Nl2BrString,
			},
		}))
		sessions := session.New(session.Options{
			SessionIdName: core.SessionName,
		})
		core.Server.Use(xsrf.New(time.Hour))
		core.Server.Use(binding.Bind())
		core.Server.Use(sessions)
		core.Server.Use(flash.Flashes(sessions))
	}
	return nil, nil
}

func (bs *BootstrapService) Install(_ interface{}) (*Result, error) {
	// create tables
	if err := core.Db.Sync2(new(model.User),
		new(model.UserToken),
		new(model.Theme),
		new(model.Article),
		new(model.ArticleTag),
		new(model.Setting),
		new(model.Media),
		new(model.Page),
		new(model.Comment)); err != nil {
		return nil, err
	}

	// insert default user
	user := &model.User{
		Name:      "admin",
		Email:     "admin@admin.com",
		Nick:      "admin",
		Profile:   "this is administrator",
		Role:      model.USER_ROLE_ADMIN,
		Status:    model.USER_STATUS_ACTIVE,
		AvatarUrl: utils.Gravatar("admin@admin.com"),
	}
	user.SetPassword("123456789")
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
		Title:       "PUGO",
		SubTitle:    "Simple Blog Engine",
		Keyword:     "pugo,blog,go,golang",
		Description: "PUGO is a simple blog engine by golang",
		HostName:    "http://localhost/",
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

	// assign install time to config
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
	}
	return nil, nil
}
