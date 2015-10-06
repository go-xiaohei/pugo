package service

import (
	"errors"
	"github.com/fuxiaohei/pugo/src/core"
	"github.com/fuxiaohei/pugo/src/model"
	"net/url"
	"strings"
)

var (
	Setting *SettingService = new(SettingService)

	ErrSettingNotFound    = errors.New("setting-not-found")
	ErrSettingMenuBadData = errors.New("setting-menu-bad-data")
)

type SettingService struct {
	General *model.SettingGeneral
	Media   *model.SettingMedia
	Content *model.SettingContent
	Comment *model.SettingComment
	Menu    []*model.SettingMenu
}

type SettingReadOption struct {
	Type         int
	UserId       int
	IsUseDefault bool // use default value if user's is not exist
}

func (ss *SettingService) Read(v interface{}) (*Result, error) {
	opt, ok := v.(SettingReadOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(ss.Read, opt)
	}
	userSetting := new(model.Setting)
	if _, err := core.Db.Where("type = ? AND user_id = ?", opt.Type, opt.UserId).Get(userSetting); err != nil {
		return nil, err
	}
	if userSetting.Id == 0 {
		if !opt.IsUseDefault {
			return nil, ErrSettingNotFound
		}
		if _, err := core.Db.Where("type = ? AND user_id = ?", opt.Type, 0).Get(userSetting); err != nil {
			return nil, err
		}
		if userSetting.Id == 0 {
			return nil, ErrSettingNotFound
		}
		return newResult(ss.Read, userSetting), nil
	}
	return newResult(ss.Read, userSetting), nil
}

func (ss *SettingService) Write(v interface{}) (*Result, error) {
	s, ok := v.(*model.Setting)
	if !ok {
		return nil, ErrServiceFuncNeedType(ss.Read, s)
	}
	if _, err := core.Db.Where("type = ? AND user_id = ?", s.Type, s.UserId).Update(s); err != nil {
		return nil, err
	}
	return nil, nil
}

func (ss *SettingService) CreateMenu(v interface{}) (*Result, error) {
	form, ok := v.(url.Values)
	if !ok {
		return nil, ErrServiceFuncNeedType(ss.CreateMenu, form)
	}
	if len(form["name"]) != len(form["link"]) {
		return nil, ErrSettingMenuBadData
	}
	if len(form["name"]) != len(form["title"]) {
		return nil, ErrSettingMenuBadData
	}
	if len(form["name"]) != len(form["new"]) {
		return nil, ErrSettingMenuBadData
	}
	menuSettings := make([]*model.SettingMenu, len(form["name"]))
	for i, v := range form["name"] {
		s := &model.SettingMenu{
			v, form["link"][i], form["title"][i],
			form["new"][i] == "true", strings.ToLower(v),
		}
		menuSettings[i] = s
	}
	return newResult(ss.CreateMenu, &menuSettings), nil
}
