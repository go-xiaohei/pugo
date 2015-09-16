package admin

import (
	"github.com/tango-contrib/xsrf"
	"pugo/src/middle"
	"pugo/src/model"
	"pugo/src/service"
	"strings"
)

type SettingGeneralController struct {
	xsrf.Checker

	middle.AuthorizeRequire
	middle.AdminRender
	middle.Validator
	middle.Responsor
}

func (sc *SettingGeneralController) Get() {
	sc.Title("GENERAL SETTING - PUGO")
	sc.Assign("XsrfHTML", sc.XsrfFormHtml())
	sc.Assign("GeneralSetting", service.Setting.General)
	sc.Assign("MediaSetting", service.Setting.Media)
	sc.Render("setting_general.tmpl")
}

type SettingGeneralForm struct {
	Title       string `form:"title"`
	SubTitle    string `form:"subtitle"`
	Keyword     string `form:"keyword"`
	Description string `form:"desc"`
	HostName    string `form:"host"`
}

func (f SettingGeneralForm) toSettingGeneral() *model.SettingGeneral {
	return &model.SettingGeneral{
		Title:       f.Title,
		SubTitle:    f.SubTitle,
		Keyword:     f.Keyword,
		Description: f.Description,
		HostName:    f.HostName,
	}
}

func (sc *SettingGeneralController) Post() {
	form := new(SettingGeneralForm)
	if err := sc.Validator.Validate(form, sc); err != nil {
		sc.JSONError(200, err)
		return
	}
	setting := &model.Setting{
		Name:   "general",
		UserId: 0,
		Type:   model.SETTING_TYPE_GENERAL,
	}
	setting.Encode(form.toSettingGeneral())
	if err := service.Call(service.Setting.Write, setting); err != nil {
		sc.JSONError(200, err)
		return
	}
	service.Setting.General = form.toSettingGeneral()
	sc.JSON(nil)
}

type SettingMediaForm struct {
	MaxFileSize int64  `form:"size"`
	ImageFile   string `form:"image"`
	DocFile     string `form:"doc"`
	CommonFile  string `form:"common"`
}

func (f SettingMediaForm) toSettingMedia() *model.SettingMedia {
	return &model.SettingMedia{
		MaxFileSize: f.MaxFileSize,
		ImageFile:   strings.Split(f.ImageFile, " "),
		DocFile:     strings.Split(f.DocFile, " "),
		CommonFile:  strings.Split(f.CommonFile, " "),
	}
}

func (sc *SettingGeneralController) PostMedia() {
	form := new(SettingMediaForm)
	if err := sc.Validator.Validate(form, sc); err != nil {
		sc.JSONError(200, err)
		return
	}
	setting := &model.Setting{
		Name:   "media",
		UserId: 0,
		Type:   model.SETTING_TYPE_MEDIA,
	}
	setting.Encode(form.toSettingMedia())
	if err := service.Call(service.Setting.Write, setting); err != nil {
		sc.JSONError(200, err)
		return
	}
	service.Setting.Media = form.toSettingMedia()
	sc.JSON(nil)
}
