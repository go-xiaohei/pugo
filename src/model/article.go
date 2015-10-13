package model

import (
	"fmt"
	"github.com/go-xiaohei/pugo/src/core"
	"github.com/go-xiaohei/pugo/src/utils"
	"html/template"
	"strings"
	"time"
)

const (
	ARTICLE_STATUS_PUBLISH = 1
	ARTICLE_STATUS_DRAFT   = 11
	ARTICLE_STATUS_DELETE  = 111

	ARTICLE_COMMENT_OPEN  = 1
	ARTICLE_COMMENT_WAIT  = 11 // waiting 30 days to close
	ARTICLE_COMMENT_CLOSE = 111

	ARTICLE_BODY_MARKDOWN = 1
	ARTICLE_BODY_HTML     = 2
)

// article struct
type Article struct {
	Id         int64
	UserId     int64  `xorm:"index(user)"`
	Title      string `xorm:"VARCHAR(100) index(title) notnull"`
	CreateTime int64  `xorm:"INT(12) created"`
	UpdateTime int64  `xorm:"INT(12) updated"`
	Link       string `xorm:"VARCHAR(100) unique notnull"`
	Preview    string `xorm:"TEXT"`
	Body       string `xorm:"TEXT notnull"`
	BodyType   int8   `xorm:"INT(8) notnull"`
	Topic      string
	TagString  string `xorm:"VARCHAR(200)"`

	Hits          int64 `xorm:"INT(8)"`
	Comments      int64 `xorm:"INT(8)"`
	Status        int8  `xorm:"INT(8)"`
	CommentStatus int8  `xorm:"INT(8)"`

	tagData   []*ArticleTag `xorm:"-"`
	userData  *User         `xorm:"-"`
	IsNewRead bool          `xorm:"-"`
}

// tag struct
type ArticleTag struct {
	Id        int64
	ArticleId int64
	Tag       string `xorm:"VARCHAR(50) notnull"`
}

// is article publish
func (a *Article) IsPublish() bool {
	return a.Status == ARTICLE_STATUS_PUBLISH
}

// is article draft
func (a *Article) IsDraft() bool {
	return a.Status == ARTICLE_STATUS_DRAFT
}

// article comment is closed
func (a *Article) IsCommentClosed() bool {
	return a.CommentStatus == ARTICLE_COMMENT_CLOSE
}

// article comment enable or not
func (a *Article) IsCommentable(duration int64) bool {
	if a.CommentStatus == ARTICLE_COMMENT_OPEN {
		return true
	}
	if a.CommentStatus == ARTICLE_COMMENT_WAIT {
		// open comment in 30 days
		if time.Now().Unix()-a.CreateTime <= 3600*24*duration {
			return true
		}
	}
	return false
}

// read article's owner
func (a *Article) User() *User {
	if a.userData == nil {
		u, err := getArticleUser(a.UserId)
		if err != nil || u == nil {
			a.userData = &User{
				Name: "Unknown",
				Nick: "Unknown",
			}
		} else {
			a.userData = u
		}
	}
	return a.userData
}

// create link for article
func (a *Article) Href() string {
	if a.IsDraft() {
		return "#"
	}
	return fmt.Sprintf("/article/%d/%s.html", a.Id, a.Link)
}

func (a *Article) PreviewContent() template.HTML {
	if a.BodyType == PAGE_BODY_MARKDOWN {
		return utils.Markdown2HTML(a.Preview)
	}
	return template.HTML(a.Preview)
}

func (a *Article) Content() template.HTML {
	if a.BodyType == PAGE_BODY_MARKDOWN {
		body := strings.Replace(a.Body, "<!--more-->", "\n", -1)
		body = strings.Replace(body, "[more]", "\n", -1)
		return utils.Markdown2HTML(body)
	}
	return template.HTML(a.Body)
}

func getArticleUser(id int64) (*User, error) {
	u := new(User)
	if _, err := core.Db.Where("id = ?", id).Get(u); err != nil {
		return nil, err
	}
	if u.Id != id {
		return nil, nil
	}
	return u, nil
}

type ArticleArchive struct {
	Id         int64
	Title      string
	Link       string
	CreateTime int64
	IsNewYear  bool `xorm:"-"`
}

func (aa *ArticleArchive) TableName() string {
	return "article"
}

func (aa *ArticleArchive) Href() string {
	return fmt.Sprintf("/article/%d/%s.html", aa.Id, aa.Link)
}
