package service

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/lunny/tango"
	"github.com/tango-contrib/binding"
	"github.com/tango-contrib/flash"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/session"
	"github.com/tango-contrib/xsrf"
	"path"
	"pugo/src/core"
	"pugo/src/model"
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
			/*Funcs: template.FuncMap{
			    "Str2HTML":             utils.Str2HTML,
			    "TimeUnixFormat":       utils.TimeUnixFormat,
			    "TimeUnixFormatFriend": utils.FriendTimeUnixFormat,
			    "FriendBytesSize":      utils.FriendBytesSize,
			},*/
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
		new(model.Theme)); err != nil {
		return nil, err
	}

	// insert default user
	user := &model.User{
		Name:    "admin",
		Email:   "admin@admin.com",
		Nick:    "admin",
		Profile: "this is administrator",
		Role:    model.USER_ROLE_ADMIN,
		Status:  model.USER_STATUS_ACTIVE,
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

	// assign install time to config
	core.Cfg.Install = fmt.Sprint(time.Now().Unix())
	if err := core.Cfg.WriteToFile(core.ConfigFile); err != nil {
		return nil, err
	}
	return nil, nil
}

type BootstrapOption struct {
	Themes bool // load themes
	I18n   bool // load languages
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
	return nil, nil
}