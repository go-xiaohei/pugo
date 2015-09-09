package core

import (
	"fmt"
	"github.com/Unknwon/com"
	"gopkg.in/ini.v1"
	"strconv"
	"time"
)

type Config struct {
	Name    string
	Version string
	Date    string
	Http    configHttp
	Db      configDb
	Created string
	Install string
}

type configHttp struct {
	Host   string
	Port   string
	Domain string
}

type configDb struct {
	Driver string
	DSN    string
}

func NewConfig() *Config {
	return &Config{
		Name:    PUGO_NAME,
		Version: PUGO_VERSION,
		Date:    PUGO_VERSION_DATE,
		Http: configHttp{
			Host:   "0.0.0.0",
			Port:   "9899",
			Domain: "localhost",
		},
		Db: configDb{
			Driver: "mysql",
			DSN:    "root:fuxiaohei@tcp/pugo?charset=utf8",
		},
		Created: fmt.Sprint(time.Now().Unix()),
		Install: "0",
	}
}

func (c *Config) Sync(file string) error {
	if !com.IsFile(file) {
		return c.WriteToFile(file)
	}
	if err := c.ReadFromFile(file); err != nil {
		return err
	}
	t, _ := time.Parse("20060102", PUGO_VERSION_DATE)
	// version date is over config file created time
	// it means the version is upgraded, config data may changes
	created, _ := strconv.ParseInt(c.Created, 10, 64)
	if t.Unix() > created {
		return c.WriteToFile(file)
	}
	return nil
}

func (c *Config) WriteToFile(file string) error {
	iniFile := ini.Empty()
	if err := iniFile.ReflectFrom(c); err != nil {
		return err
	}
	return iniFile.SaveToIndent(file, "  ")
}

func (c *Config) ReadFromFile(file string) error {
	iniFile, err := ini.Load(file)
	if err != nil {
		return err
	}
	return iniFile.MapTo(c)
}
