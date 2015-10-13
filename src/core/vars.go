package core

import (
	"github.com/codegangsta/cli"
	"github.com/go-xorm/xorm"
	"github.com/lunny/tango"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	PUGO_VERSION       = "1.0"
	PUGO_VERSION_STATE = "alpha"
	PUGO_VERSION_DATE  = "20150910"
	PUGO_NAME          = "Pugo"
	PUGO_DESCRIPTION   = "a simple golang blog engine"
	PUGO_AUTHOR        = "fuxiaohei"
	PUGO_AUTHOR_EMAIL  = "fuxiaohei@vip.qq.com"

	RUN_MOD_PROD  = "prod"
	RUN_MOD_DEBUG = "debug"
)

var (
	RunMode = RUN_MOD_PROD // prod || debug

	App    *cli.App
	Cfg    *Config
	Db     *xorm.Engine
	Server *tango.Tango
	Crash  log15.Logger

	ConfigFile      string = "config.ini"
	StaticPrefix    string = "/static"
	StaticDirectory string = "static"
	ThemePrefix     string = "/theme"
	ThemeDirectory  string = "theme"
	SessionName     string = "PUGO_SESSION"
	CrashLogFile    string = "crash.log"
	BackupDirectory string = "backup"
)
