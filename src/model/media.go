package model

const (
	MEDIA_TYPE_IMAGE = iota + 1
	MEDIA_TYPE_DOC
	MEDIA_TYPE_COMMON
)

type Media struct {
	Id         int64
	UserId     int64
	Name       string `xorm:"VARCHAR(255)"`
	FileName   string `xorm:"VARCHAR(255)"`
	FilePath   string `xorm:"not null"`
	FileSize   int64  `xorm:"INT(12)"`
	FileType   int    `xorm:"INT(8) notnull index(type)"`
	CreateTime int64  `xorm:"INT(12) created"`
	Downloads  int    `xorm:"INT(8)"`
}

func (m *Media) Type() string {
	if m.FileType == MEDIA_TYPE_IMAGE {
		return "image"
	}
	if m.FileType == MEDIA_TYPE_DOC {
		return "doc"
	}
	return "common"
}
