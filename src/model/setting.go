package model

import (
	"encoding/json"
	"strings"
)

const (
	SETTING_TYPE_GENERAL = iota + 1
	SETTING_TYPE_MEDIA
)

type Setting struct {
	Id         int64
	Name       string `xorm:"VARCHAR(50) notnull index(name)"`
	Value      string `xorm:"TEXT notnull"`
	UserId     int64
	Type       int8  `xorm:"INT(8) index(type)"`
	CreateTime int64 `xorm:"INT(12) created"`
}

func (s *Setting) Encode(v interface{}) {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	s.Value = string(bytes)
}

func (s *Setting) ToGeneral() *SettingGeneral {
	general := new(SettingGeneral)
	if err := json.Unmarshal([]byte(s.Value), general); err != nil {
		panic(err)
	}
	return general
}

func (s *Setting) ToMedia() *SettingMedia {
	media := new(SettingMedia)
	if err := json.Unmarshal([]byte(s.Value), media); err != nil {
		panic(err)
	}
	return media
}

type SettingGeneral struct {
	Title       string `json:"title"`
	SubTitle    string `json:"sub_title"`
	Keyword     string `json:"keyword_meta"`
	Description string `json:"description_meta"`
	HostName    string `json:"host_name"`
}

type SettingMedia struct {
	MaxFileSize int64    `json:"max_file_size"`
	ImageFile   []string `json:"image_file"`
	DocFile     []string `json:"doc_file"`
	CommonFile  []string `json:"common_file"`
}

func (sm SettingMedia) Image() string {
	return strings.Join(sm.ImageFile, " ")
}

func (sm SettingMedia) Doc() string {
	return strings.Join(sm.DocFile, " ")
}

func (sm SettingMedia) Common() string {
	return strings.Join(sm.CommonFile, " ")
}