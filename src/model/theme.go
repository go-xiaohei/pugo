package model

import "github.com/Unknwon/com"

const (
	THEME_STATUS_LOCKED = iota + 1
	THEME_STATUS_INVALID
	THEME_STATUS_NORMAL
	THEME_STATUS_CURRENT
)

type Theme struct {
	Id          int64
	Name        string `xorm:"VARCHAR(100) unique"`
	Author      string `xorm:"VARCHAR(100) notnull"`
	AuthorUrl   string `xorm:"VARCHAR(200)"`
	Version     string `xorm:"VARCHAR(20) notnull"`
	Directory   string `xorm:"VARCHAR(200) notnull"`
	InstallTime int64  `xorm:"INT(12) created"`
	Status      int    `xorm:"INT(4)"`
}

func (t *Theme) IsLocked() bool {
	return t.Status == THEME_STATUS_LOCKED // locked theme
}

func (t *Theme) IsCurrent() bool {
	return t.Status == THEME_STATUS_CURRENT // current theme
}

func (t *Theme) IsValid() bool {
	return com.IsDir(t.Directory)
}
