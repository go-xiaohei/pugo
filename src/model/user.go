package model

import "golang.org/x/crypto/bcrypt"

const (
	USER_ROLE_NORMAL = 2 << iota // normal user
	USER_ROLE_ADMIN              //admin user
)
const (
	USER_STATUS_ACTIVE  = 2 << iota // active user
	USER_STATUS_FRESH               // fresh user, means need activate
	USER_STATUS_BLOCKED             // blocked user
	USER_STATUS_REMOVED             // removed user
)

type User struct {
	Id       int64
	Name     string `xorm:"VARCHAR(200) notnull unique"`
	Password string `xorm:"VARCHAR(100) notnull"`
	Email    string `xorm:"VARCHAR(200) notnull unique"`

	Nick      string `xorm:"VARCHAR(50)"`
	Profile   string
	Url       string `xorm:"VARCHAR(200)"`
	AvatarUrl string `xorm:"VARCHAR(200)"`

	CreateTime    int64 `xorm:"INT(12) created"`
	UpdateTime    int64 `xorm:"INT(12) updated"`
	LastLoginTime int64 `xorm:"INT(12)"`

	Role   int8 `xorm:"INTEGER(8) index(role)"`
	Status int8 `xorm:"INTEGER(8) index(status)"`
}

func (u *User) IsAdmin() bool {
	return u.Role == USER_ROLE_ADMIN
}

func (u *User) IsAccessible() bool {
	return u.Status <= USER_STATUS_FRESH
}

func (u *User) IsPassword(raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(raw)) == nil
}

func (u *User) SetPassword(pwd string) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	u.Password = string(bytes)
}

type FrontUser struct {
	Id        int64
	Name      string
	Email     string
	Nick      string
	AvatarUrl string
	Role      int8
}

func NewFrontUser(u *User) *FrontUser {
	return &FrontUser{
		u.Id, u.Name, u.Email, u.Nick, u.AvatarUrl, u.Role,
	}
}
