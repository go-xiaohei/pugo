package model

import (
	"github.com/fuxiaohei/pugo/src/core"
	"github.com/fuxiaohei/pugo/src/utils"
)

const (
	COMMENT_FROM_ARTICLE = iota + 1
	COMMENT_FROM_PAGE
)

const (
	COMMENT_STATUS_APPROVED = iota + 1
	COMMENT_STATUS_WAIT
	COMMENT_STATUS_SPAM
	COMMENT_STATUS_DELETED
)

type Comment struct {
	Id         int64  `json:"id"`
	Name       string `xorm:"VARCHAR(100) notnull" json:"name"`
	UserId     int64  `json:"user_id"`
	Email      string `xorm:"VARCHAR(200) notnull" json:"-"`
	Url        string `xorm:"VARCHAR(200)" json:"url"`
	AvatarUrl  string `xorm:"VARCHAR(200)" json:"avatar"`
	Body       string `xorm:"TEXT notnull" json:"body"`
	CreateTime int64  `xorm:"created" json:"created"`
	Status     int    `xorm:"INT(8) index(status)" json:"status"`

	UserIp    string `xorm:"VARCHAR(200)" json:"ip"`
	UserAgent string `xorm:"VARCHAR(200)" json:"user_agent"`

	From     int   `xorm:"INT(8) index(from)" json:"-"`
	FromId   int64 `xorm:"index(from)" json:"-"`
	ParentId int64 `xorm:"index(parent)" json:"parent"`

	parent *Comment `xorm:"-"`
}

func (c *Comment) IsTopApproved() bool {
	return c.Status == COMMENT_STATUS_APPROVED && c.ParentId == 0
}

func (c *Comment) AuthorUrl() string {
	if c.Url == "" {
		return "#"
	}
	return c.Url
}

func (c *Comment) IsApproved() bool {
	return c.Status == COMMENT_STATUS_APPROVED
}

func (c *Comment) IsWait() bool {
	return c.Status == COMMENT_STATUS_WAIT
}

func (c *Comment) IsSpam() bool {
	return c.Status == COMMENT_STATUS_SPAM
}

func (c *Comment) FromTitle() string {
	if c.From == COMMENT_FROM_ARTICLE {
		if article := getArticleById(c.FromId); article != nil {
			return article.Title
		}
	}
	if c.From == COMMENT_FROM_PAGE {
		if page := getPageById(c.FromId); page != nil {
			return page.Title
		}
	}
	return ""
}

func (c *Comment) GetParent() *Comment {
	if c.ParentId == 0 {
		return nil
	}
	co := new(Comment)
	if _, err := core.Db.Where("id = ?", c.ParentId).Get(co); err != nil {
		return nil
	}
	if c.ParentId != co.ParentId {
		return nil
	}
	return co
}

func getArticleById(id int64) *Article {
	a := new(Article)
	if _, err := core.Db.Where("id = ?", id).Get(a); err != nil {
		return nil
	}
	if a.Id != id {
		return nil
	}
	return a
}

func getPageById(id int64) *Page {
	a := new(Page)
	if _, err := core.Db.Where("id = ?", id).Get(a); err != nil {
		return nil
	}
	if a.Id != id {
		return nil
	}
	return a
}

type FrontComment struct {
	Id         int64  `json:"id"`
	Name       string `xorm:"VARCHAR(100) notnull" json:"name"`
	UserId     int64  `json:"user_id"`
	Url        string `xorm:"VARCHAR(200)" json:"url"`
	AvatarUrl  string `xorm:"VARCHAR(200)" json:"avatar"`
	Body       string `xorm:"TEXT notnull" json:"body"`
	CreateTime string `xorm:"created" json:"created"`
	Status     int    `xorm:"INT(8) index(status)" json:"status"`

	UserIp    string `xorm:"VARCHAR(200)" json:"ip"`
	UserAgent string `xorm:"VARCHAR(200)" json:"user_agent"`
	ParentId  int64  `xorm:"index(parent)" json:"parent"`
}

func NewFrontComment(c *Comment) *FrontComment {
	fc := &FrontComment{
		Id:         c.Id,
		Name:       c.Name,
		UserId:     c.UserId,
		Url:        c.Url,
		AvatarUrl:  c.AvatarUrl,
		Body:       utils.Nl2BrString(c.Body),
		CreateTime: utils.TimeUnixFriend(c.CreateTime),
		Status:     c.Status,
		UserIp:     c.UserIp,
		UserAgent:  c.UserAgent,
		ParentId:   c.ParentId,
	}
	return fc
}
