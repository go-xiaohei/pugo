package model

import "html/template"

const (
	MESSAGE_FROM_ARTICLE = iota + 1
	MESSAGE_FROM_PAGE
	MESSAGE_FROM_COMMENT
	MESSAGE_FROM_MEDIA
	MESSAGE_FROM_BACKUP
)

const (
	MESSAGE_TYPE_ARTICLE_CREATE   = 101
	MESSAGE_TYPE_ARTICLE_UPDATE   = 102
	MESSAGE_TYPE_ARTICLE_REMOVE   = 103
	MESSAGE_TYPE_PAGE_CREATE      = 201
	MESSAGE_TYPE_PAGE_UPDATE      = 202
	MESSAGE_TYPE_PAGE_REMOVE      = 203
	MESSAGE_TYPE_COMMENT_CREATE   = 301
	MESSAGE_TYPE_COMMENT_REMOVE   = 302
	MESSAGE_TYPE_COMMENT_REPLY    = 303
	MESSAGE_TYPE_COMMENT_FEEDBACK = 304 // means reply from admin panel
	MESSAGE_TYPE_MEDIA_UPLOAD     = 401
	MESSAGE_TYPE_BACKUP_CREATE    = 501
)

type Message struct {
	Id         int64
	UserId     int64  `xorm:"notnull"`
	From       int    `xorm:"INT(8) notnull"`
	FromId     int64  `xorm:"notnull"`
	Type       int    `xorm:"INT(8) notnull"`
	Body       string `xorm:"TEXT"`
	CreateTime int64  `xorm:"INT(12) created"`
}

func (m *Message) HTML() template.HTML {
	return template.HTML(m.Body)
}
