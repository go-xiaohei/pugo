package service

import (
	"errors"
	"fmt"
	"github.com/go-xiaohei/pugo/src/core"
	"github.com/go-xiaohei/pugo/src/model"
	"time"
)

var (
	User *UserService = new(UserService)

	ErrUserNotFound      = errors.New("user-not-found")
	ErrUserNotAccess     = errors.New("user-not-access")
	ErrUserWrongPassword = errors.New("user-wrong-password")
	ErrTokenNotFound     = errors.New("token-not-found")
	ErrTokenExpired      = errors.New("token-expired")
	ErrPasswordConfirm   = errors.New("password-confirm-error")
)

type UserService struct{}

type UserAuthOption struct {
	Name           string // auth by name
	Email          string // auth by email
	Password       string
	ExpireDuration int64
	Origin         string
}

// authorize user
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
		if user, err = getUserBy("name", opt.Name); err != nil {
			return nil, err
		}
	} else if opt.Email != "" {
		if user, err = getUserBy("email", opt.Email); err != nil {
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

func getUserBy(col string, value interface{}) (*model.User, error) {
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

type UserVerifyOption struct {
	Hash           string
	Origin         string
	Extend         bool
	ExtendDuration int64
}

func (us *UserService) Verify(v interface{}) (*Result, error) {
	opt, ok := v.(UserVerifyOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(us.Verify, opt)
	}

	// get token
	token, err := us.getToken(opt.Hash, opt.Origin)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, ErrTokenNotFound
	}
	if token.IsExpired() {
		return nil, ErrTokenExpired
	}

	// get token's owner
	user, err := getUserBy("id", token.UserId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if !user.IsAccessible() {
		return nil, ErrUserNotAccess
	}

	// extend user
	if opt.Extend {
		token.ExpireTime += opt.ExtendDuration
		if err := us.extendToken(token.Id, token.ExpireTime); err != nil {
			return nil, err
		}
	}
	res := newResult(us.Verify)
	res.Set(user, token)
	return res, nil
}

func (us *UserService) getToken(hash, origin string) (*model.UserToken, error) {
	t := new(model.UserToken)
	if _, err := core.Db.Where("`from` = ? AND hash = ?", origin, hash).Get(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (us *UserService) extendToken(id, newExpire int64) error {
	if _, err := core.Db.Exec("UPDATE user_token SET expire_time = ? WHERE id = ?", newExpire, id); err != nil {
		return err
	}
	return nil
}

func (us *UserService) Unauthorize(v interface{}) (*Result, error) {
	token, ok := v.(*string)
	if !ok {
		return nil, ErrServiceFuncNeedType(us.Unauthorize, token)
	}
	if _, err := core.Db.Where("hash = ?", *token).Delete(new(model.UserToken)); err != nil {
		return nil, err
	}
	return nil, nil
}

func (us *UserService) Profile(v interface{}) (*Result, error) {
	u, ok := v.(*model.User)
	if !ok {
		return nil, ErrServiceFuncNeedType(us.Profile, u)
	}
	if _, err := core.Db.Where("id = ?", u.Id).Cols("name,email,profile,nick,url").Update(u); err != nil {
		return nil, err
	}
	return nil, nil
}

type UserPasswordOption struct {
	Old     string
	New     string
	Confirm string
	User    *model.User
}

func (us *UserService) Password(v interface{}) (*Result, error) {
	opt, ok := v.(UserPasswordOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(us.Password, opt)
	}
	if opt.New != opt.Confirm {
		return nil, ErrPasswordConfirm
	}
	if !opt.User.IsPassword(opt.Old) {
		return nil, ErrUserWrongPassword
	}
	oldPwd := opt.User.Password
	opt.User.SetPassword(opt.New)
	if _, err := core.Db.Where("id = ?", opt.User.Id).Cols("password").Update(opt.User); err != nil {
		opt.User.Password = oldPwd
		return nil, err
	}
	return nil, nil
}
