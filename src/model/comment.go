package model

import (
	"github.com/go-xiaohei/pugo/src/core"
	"github.com/go-xiaohei/pugo/src/utils"
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
	CreateTime int64  `xorm:"INT(12) created" json:"created"`
	Status     int    `xorm:"INT(8) index(status)" json:"status"`

	UserIp    string `xorm:"VARCHAR(200)" json:"ip"`
	UserAgent string `xorm:"VARCHAR(200)" json:"user_agent"`

	From     int   `xorm:"INT(8) index(from)" json:"-"`
	FromId   int64 `xorm:"index(from)" json:"-"`
	ParentId int64 `xorm:"index(parent)" json:"parent"`

	parent    *Comment `xorm:"-"`
	FromTitle string   `xorm:"-"`
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

func (c *Comment) GetParent() *Comment {
	if c.ParentId == 0 {
		return nil
	}
	if c.parent == nil {
		co := new(Comment)
		if _, err := core.Db.Where("id = ?", c.ParentId).Get(co); err != nil {
			return nil
		}
		if c.ParentId != co.Id {
			return nil
		}
		c.parent = co
	}
	return c.parent
}

func (c *Comment) GetTitle() string {
	if c.FromTitle == "" {
		if c.From == COMMENT_FROM_ARTICLE {
			c.FromTitle = getArticleTitleById(c.FromId)
		}
		if c.From == COMMENT_FROM_PAGE {
			c.FromTitle = getPageTitleById(c.FromId)
		}
	}
	return c.FromTitle
}

func getArticleTitleById(id int64) string {
	a := new(Article)
	if _, err := core.Db.Cols("id,title").Where("id = ?", id).Get(a); err != nil {
		return ""
	}
	if a.Id != id {
		return ""
	}
	return a.Title
}

func getPageTitleById(id int64) string {
	a := new(Page)
	if _, err := core.Db.Cols("id,title").Where("id = ?", id).Get(a); err != nil {
		return ""
	}
	if a.Id != id {
		return ""
	}
	return a.Title
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

type CommentsGroup struct {
	comments          []*Comment
	cacheParent       map[int64]*Comment
	cacheArticleTitle map[int64]string
	cachePageTitle    map[int64]string
}

func NewCommentsGroup(cmts []*Comment) *CommentsGroup {
	return &CommentsGroup{
		comments:          cmts,
		cacheParent:       make(map[int64]*Comment),
		cacheArticleTitle: make(map[int64]string),
		cachePageTitle:    make(map[int64]string),
	}
}

func (cg *CommentsGroup) FillAll() {
	for _, c := range cg.comments {
		if c.From == COMMENT_FROM_ARTICLE {
			if cg.cacheArticleTitle[c.FromId] == "" {
				cg.cacheArticleTitle[c.FromId] = getArticleTitleById(c.FromId)
			}
			c.FromTitle = cg.cacheArticleTitle[c.FromId]
		}
		if c.From == COMMENT_FROM_PAGE {
			if cg.cachePageTitle[c.FromId] == "" {
				cg.cachePageTitle[c.FromId] = getPageTitleById(c.FromId)
			}
			c.FromTitle = cg.cachePageTitle[c.FromId]
		}
		if c.ParentId > 0 {
			if cg.cacheParent[c.ParentId] == nil {
				co := new(Comment)
				if _, err := core.Db.Where("id = ?", c.ParentId).Get(co); err != nil {
					continue
				}
				if c.ParentId != co.Id {
					continue
				}
				cg.cacheParent[c.ParentId] = co
			}
			c.parent = cg.cacheParent[c.ParentId]
		}
	}
}

func (cg *CommentsGroup) Comments() []*Comment {
	return cg.comments
}
