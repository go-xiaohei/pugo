package model

import (
	"fmt"
	"github.com/go-xiaohei/pugo/src/utils"
	"strings"
	"time"
)

const (
	PAGE_STATUS_PUBLISH = 1
	PAGE_STATUS_DRAFT   = 11
	PAGE_STATUS_DELETE  = 111

	PAGE_COMMENT_OPEN  = 1
	PAGE_COMMENT_WAIT  = 11 // waiting 30 days to close
	PAGE_COMMENT_CLOSE = 111

	PAGE_BODY_MARKDOWN = 1
	PAGE_BODY_HTML     = 2
)

type Page struct {
	Id         int64
	UserId     int64  `xorm:"index(user)"`
	Title      string `xorm:"VARCHAR(100) index(title) notnull"`
	CreateTime int64  `xorm:"INT(12) created"`
	UpdateTime int64  `xorm:"INT(12) updated"`
	Link       string `xorm:"VARCHAR(100) unique notnull"`
	TopLink    bool   `xorm:"INT(1)"`
	Body       string `xorm:"TEXT notnull"`
	BodyType   int8   `xorm:"INT(8) notnull"`
	Template   string `xorm:"VARCHAR(50)"`

	Hits          int64 `xorm:"INT(8)"`
	Comments      int64 `xorm:"INT(8)"`
	Status        int8  `xorm:"INT(8)"`
	CommentStatus int8  `xorm:"INT(8)"`

	userData *User `xorm:"-"`
}

func (p *Page) Href() string {
	if p.IsDraft() {
		return "#"
	}
	if p.TopLink {
		return "/" + p.Link + ".html"
	}
	return fmt.Sprintf("/page/%d/%s.html", p.Id, p.Link)
}

// is article publish
func (p *Page) IsPublish() bool {
	return p.Status == ARTICLE_STATUS_PUBLISH
}

// is article draft
func (p *Page) IsDraft() bool {
	return p.Status == ARTICLE_STATUS_DRAFT
}

// article comment is closed
func (p *Page) IsCommentClosed() bool {
	return p.CommentStatus == ARTICLE_COMMENT_CLOSE
}

// article comment enable or not
func (p *Page) IsCommentable() bool {
	if p.CommentStatus == ARTICLE_COMMENT_OPEN {
		return true
	}
	if p.CommentStatus == ARTICLE_COMMENT_WAIT {
		// open comment in 30 days
		if time.Now().Unix()-p.CreateTime <= 3600*24*30 {
			return true
		}
	}
	return false
}

// read article's owner
func (p *Page) User() *User {
	if p.userData == nil {
		u, err := getArticleUser(p.UserId)
		if err != nil || u == nil {
			p.userData = &User{
				Name: "Unknown",
				Nick: "Unknown",
			}
		} else {
			p.userData = u
		}
	}
	return p.userData
}

func (p *Page) Content() string {
	if p.BodyType == PAGE_BODY_MARKDOWN {
		body := strings.Replace(p.Body, "<!--more-->", "\n", -1)
		body = strings.Replace(body, "[more]", "\n", -1)
		return utils.Markdown2String(body)
	}
	return p.Body
}
