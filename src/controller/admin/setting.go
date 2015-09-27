package admin

import (
	"github.com/fuxiaohei/pugo/src/middle"
	"github.com/fuxiaohei/pugo/src/model"
	"github.com/fuxiaohei/pugo/src/service"
	"github.com/tango-contrib/xsrf"
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
	DynamicLink string `form:"dync"`
}

func (f SettingMediaForm) toSettingMedia() *model.SettingMedia {
	return &model.SettingMedia{
		MaxFileSize: f.MaxFileSize,
		ImageFile:   strings.Split(f.ImageFile, " "),
		DocFile:     strings.Split(f.DocFile, " "),
		CommonFile:  strings.Split(f.CommonFile, " "),
		DynamicLink: f.DynamicLink == "true",
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

type SettingThemeController struct {
	middle.AuthorizeRequire
	middle.AdminRender
}

func (stc *SettingThemeController) Get() {
	themes := make([]*model.Theme, 0)
	if err := service.Call(service.Theme.All, nil, &themes); err != nil {
		stc.RenderError(500, err)
		return
	}
	stc.Title("THEME - PUGO")
	stc.Assign("Themes", themes)
	stc.Render("setting_theme.tmpl")
}

type SettingContentForm struct {
	PageSize         int    `form:"page_size"`
	RSSFullText      bool   `form:"rss_full_text"`
	RSSNumberLimit   int    `form:"rss_number"`
	TopPage          int64  `form:"top_page"`
	PageDisallowLink string `form:"disallow_link"`
}

func (scf SettingContentForm) toSettingContent() *model.SettingContent {
	return &model.SettingContent{
		PageSize:         scf.PageSize,
		RSSFullText:      scf.RSSFullText,
		RSSNumberLimit:   scf.RSSNumberLimit,
		TopPage:          scf.TopPage,
		PageDisallowLink: strings.Split(scf.PageDisallowLink, " "),
	}
}

type SettingContentController struct {
	xsrf.Checker

	middle.AuthorizeRequire
	middle.AdminRender
	middle.Validator
	middle.Responsor
}

func (sc *SettingContentController) Get() {
	sc.Title("GENERAL CONTENT - PUGO")
	sc.Assign("XsrfHTML", sc.XsrfFormHtml())
	sc.Assign("ContentSetting", service.Setting.Content)
	sc.Render("setting_content.tmpl")
}

func (sc *SettingContentController) Post() {
	form := new(SettingContentForm)
	if err := sc.Validator.Validate(form, sc); err != nil {
		sc.JSONError(200, err)
		return
	}
	setting := &model.Setting{
		Name:   "content",
		UserId: 0,
		Type:   model.SETTING_TYPE_CONTENT,
	}
	setting.Encode(form.toSettingContent())
	if err := service.Call(service.Setting.Write, setting); err != nil {
		sc.JSONError(200, err)
		return
	}
	service.Setting.Content = form.toSettingContent()
	sc.JSON(nil)
}
