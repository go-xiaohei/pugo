package service

import (
	"errors"
	"pugo/src/core"
	"pugo/src/model"
)

var (
	Setting *SettingService = new(SettingService)

	ErrSettingNotFound = errors.New("setting-not-found")
)

type SettingService struct {
	General *model.SettingGeneral
	Media   *model.SettingMedia
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