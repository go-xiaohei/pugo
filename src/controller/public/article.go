package public

import (
	"github.com/lunny/tango"
	"pugo/src/middle"
	"pugo/src/model"
	"pugo/src/service"
)

type ArticleController struct {
	tango.Ctx
	middle.AuthorizeCheck
	middle.ThemeRender
}

func (ac *ArticleController) Get() {
	var (
		article = new(model.Article)
		opt     = service.ArticleReadOption{
			Id:     ac.ParamInt64(":id"),
			Link:   ac.Param(":link"),
			Status: model.ARTICLE_STATUS_PUBLISH,
		}
	)
	if err := service.Call(service.Article.Read, opt, article); err != nil {
		ac.RenderError(500, err)
		return
	}
	if article.Id != opt.Id || article.Link != opt.Link {
		ac.RenderError(404, nil)
		return
	}
	ac.Title(article.Title + " - " + service.Setting.General.Title)
	ac.Assign("Article", article)
	ac.Render("single.tmpl")
}
