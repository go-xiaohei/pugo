package model

import (
	"pugo/src/core"
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

	tagData  []*ArticleTag `xorm:"-"`
	userData *User         `xorm:"-"`
}

// tag struct
type ArticleTag struct {
	Id        int64
	ArticleId int64
	Tag       string `xorm:"VARCHAR(50) notnull"`
}

// article comment enable or not
func (a *Article) IsCommentable() bool {
	if a.CommentStatus == ARTICLE_COMMENT_OPEN {
		return true
	}
	if a.CommentStatus == ARTICLE_COMMENT_WAIT {
		// open comment in 30 days
		if time.Now().Unix()-a.CreateTime <= 3600*24*30 {
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
