package service

import (
	"errors"
	"fmt"
	"pugo/src/core"
	"pugo/src/model"
	"time"
)

var (
	User *UserService = new(UserService)

	ErrUserNotFound      = errors.New("user-not-found")
	ErrUserNotAccess     = errors.New("user-not-access")
	ErrUserWrongPassword = errors.New("user-wrong-password")
)

type UserService struct{}

type UserAuthOption struct {
	Name           string // auth by name
	Email          string // auth by email
	Password       string
	ExpireDuration int64
	Origin         string
}

func (us *UserService) Authorize(v interface{}) (*Result, error) {
	opt, ok := v.(UserAuthOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(us.Authorize, opt)
	}
	var (
		user *model.User
		err  error
	)

	// get user
	if opt.Name != "" {
		if user, err = us.getUserBy("name", opt.Name); err != nil {
			return nil, err
		}
	} else if opt.Email != "" {
		if user, err = us.getUserBy("email", opt.Email); err != nil {
			return nil, err
		}
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if !user.IsAccessible() {
		return nil, ErrUserNotAccess
	}
	if !user.IsPassword(opt.Password) {
		return nil, ErrUserWrongPassword
	}

	// create token
	var token *model.UserToken
	if token, err = us.createToken(user, opt); err != nil {
		return nil, err
	}

	// update login time
	if err = us.updateLoginTime(user); err != nil {
		return nil, err
	}

	res := newResult(us.Authorize)
	res.Set(user, token)
	return res, nil
}

// get user by column and value
func (us *UserService) getUserBy(col string, value interface{}) (*model.User, error) {
	u := new(model.User)
	if _, err := core.Db.Where(col+" = ?", value).Get(u); err != nil {
		return nil, err
	}
	if u.Id == 0 {
		return nil, nil
	}
	return u, nil
}

// create new token
func (us *UserService) createToken(u *model.User, opt UserAuthOption) (*model.UserToken, error) {
	token := &model.UserToken{
		UserId:     u.Id,
		ExpireTime: time.Now().Unix() + opt.ExpireDuration,
		From:       opt.Origin,
	}
	token.SetHash(fmt.Sprintf("%d.%d.%d", u.Id, time.Now().Unix(), opt.ExpireDuration))
	if _, err := core.Db.Insert(token); err != nil {
		return nil, err
	}
	return token, nil
}

// update login time
func (us *UserService) updateLoginTime(u *model.User) error {
	u.LastLoginTime = time.Now().Unix()
	if _, err := core.Db.Exec("UPDATE user SET last_login_time = ? WHERE id = ?", u.LastLoginTime, u.Id); err != nil {
		return err
	}
	return nil
}
