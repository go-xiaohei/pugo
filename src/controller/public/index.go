package public

import (
	"fmt"
	"github.com/lunny/tango"
	"net/http"
	"pugo/src/middle"
	"pugo/src/model"
	"pugo/src/service"
	"pugo/src/utils"
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
			Page:     1,
			Size:     5,
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
	ic.Render("index.tmpl")
}
