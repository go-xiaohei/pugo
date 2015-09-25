package model

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

	article *Article `xorm:"-"`
	page    *Page    `xorm:"-"`
	parent  *Comment `xorm:"-"`
}

func (c *Comment) IsTopApproved() bool {
	return c.Status == COMMENT_STATUS_APPROVED && c.ParentId == 0
}
