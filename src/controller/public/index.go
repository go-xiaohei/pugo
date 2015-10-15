package public

import (
	"fmt"
	"github.com/go-xiaohei/pugo/src/middle"
	"github.com/go-xiaohei/pugo/src/model"
	"github.com/go-xiaohei/pugo/src/service"
	"github.com/go-xiaohei/pugo/src/utils"
	"github.com/lunny/tango"
	"net/http"
	"time"
)

var (
	readTimeCookieName = "PUGO_READ_TIME"
)

type IndexController struct {
	tango.Ctx

	middle.ThemeRender
}

func (ic *IndexController) getReadTime() int64 {
	return ic.CookieInt64(readTimeCookieName)
}

func (ic *IndexController) setReadTime() {
	ic.Cookies().Set(&http.Cookie{
		Name:     readTimeCookieName,
		Value:    fmt.Sprint(time.Now().Unix()),
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * 10 * time.Hour),
		MaxAge:   3600 * 24 * 10 * 365,
		HttpOnly: true,
	})
}

func (ic *IndexController) Get() {
	ic.Title(service.Setting.General.FullTitle())
	var (
		opt = service.ArticleListOption{
			Status:   model.ARTICLE_STATUS_PUBLISH,
			Order:    "create_time DESC",
			Page:     ic.ParamInt(":page", 1),
			Size:     service.Setting.Content.PageSize,
			IsCount:  true,
			ReadTime: ic.getReadTime(),
		}
		articles = make([]*model.Article, 0)
		pager    = new(utils.Pager)
	)
	if err := service.Call(service.Article.List, opt, &articles, pager); err != nil {
		ic.RenderError(500, err)
		return
	}
	ic.setReadTime()
	ic.Assign("Articles", articles)
	ic.Assign("Pager", pager)
	ic.Render("index.tmpl")
}

type TagController struct {
	tango.Ctx
	middle.ThemeRender
}

func (tc *TagController) Get() {
	tag := tc.Param("tag")
	if tag == "" {
		tc.Redirect("/")
		return
	}
	tc.Title(tag + " - " + service.Setting.General.FullTitle())
	var (
		opt = service.ArticleListOption{
			Status:  model.ARTICLE_STATUS_PUBLISH,
			Order:   "create_time DESC",
			Page:    tc.ParamInt(":page", 1),
			Size:    service.Setting.Content.PageSize,
			IsCount: true,
			Tag:     tag,
		}
		articles = make([]*model.Article, 0)
		pager    = new(utils.Pager)
	)
	if err := service.Call(service.Article.ListByTag, opt, &articles, pager); err != nil {
		tc.RenderError(500, err)
		return
	}
	tc.Assign("IsTagPage", true)
	tc.Assign("Tag", tag)
	tc.Assign("TagLink", "/tag/"+tag)
	tc.Assign("Articles", articles)
	tc.Assign("Pager", pager)
	tc.Render("index.tmpl")
}
