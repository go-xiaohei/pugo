package admin

import (
	"github.com/go-xiaohei/pugo/src/middle"
	"github.com/go-xiaohei/pugo/src/model"
	"github.com/go-xiaohei/pugo/src/service"
	"github.com/go-xiaohei/pugo/src/utils"
	"github.com/lunny/tango"
	"github.com/tango-contrib/xsrf"
	"strings"
)

type ArticleWriteController struct {
	tango.Ctx
	xsrf.Checker

	middle.AuthorizeRequire
	middle.AdminRender
	middle.Validator
	middle.Responsor
}

func (awc *ArticleWriteController) Get() {
	awc.Title("WRITE ARTICLE - PUGO")
	awc.Assign("XsrfHTML", awc.XsrfFormHtml())
	if id := awc.FormInt64("id"); id > 0 {
		var (
			opt     = service.ArticleReadOption{Id: id}
			article = new(model.Article)
		)
		if err := service.Call(service.Article.Read, opt, article); err != nil {
			awc.RenderError(500, err)
			return
		}
		awc.Assign("Article", article)
	}
	awc.Render("write_article.tmpl")
}

// article post form
type ArticleForm struct {
	Title   string `form:"title" binding:"Required"`
	Link    string `form:"link" binding:"Required;AlphaDashDot"`
	Body    string `form:"body" binding:"Required"`
	Type    string `form:"type"`
	Tag     string `form:"tag"`
	Draft   string `form:"draft"`
	Id      int64  `form:"id"`
	UserId  int64
	Comment string `form:"comment"`
}

func (f ArticleForm) toArticle() *model.Article {
	article := &model.Article{
		Id:            f.Id,
		UserId:        f.UserId,
		Title:         f.Title,
		Link:          f.Link,
		Body:          f.Body,
		TagString:     f.Tag,
		Status:        model.ARTICLE_STATUS_PUBLISH,
		CommentStatus: model.ARTICLE_COMMENT_OPEN,
		Hits:          1,
	}
	if strings.Contains(f.Body, "<!--more-->") {
		article.Preview = strings.Split(f.Body, "<!--more-->")[0]
	}
	switch strings.ToLower(f.Type) {
	case "html":
		article.BodyType = model.ARTICLE_BODY_HTML
	case "markdown":
		article.BodyType = model.ARTICLE_BODY_MARKDOWN
	default:
		article.BodyType = model.ARTICLE_BODY_MARKDOWN
	}
	if f.Draft == "draft" {
		article.Status = model.ARTICLE_STATUS_DRAFT
	}
	if f.Comment == "close" {
		article.CommentStatus = model.ARTICLE_COMMENT_CLOSE
	}
	return article
}

func (awc *ArticleWriteController) Post() {
	form := new(ArticleForm)
	if err := awc.Validator.Validate(form, awc); err != nil {
		awc.JSONError(200, err)
		return
	}
	form.UserId = awc.AuthUser.Id
	var article = new(model.Article)
	if err := service.Call(service.Article.Write, form.toArticle(), article); err != nil {
		awc.JSONError(200, err)
		return
	}
	awc.JSON(map[string]interface{}{
		"article": article,
	})
}

type ArticleManageController struct {
	tango.Ctx
	middle.AuthorizeRequire
	middle.AdminRender
}

func (amc *ArticleManageController) Get() {
	amc.Title("ARTICLES - PUGO")
	var (
		opt = service.ArticleListOption{
			IsCount: true,
			Page:    amc.FormInt("page", 0),
		}
		articles = make([]*model.Article, 0)
		pager    = new(utils.Pager)
	)
	if err := service.Call(service.Article.List, opt, &articles, pager); err != nil {
		amc.RenderError(500, err)
		return
	}
	amc.Assign("Articles", articles)
	amc.Assign("Pager", pager)
	amc.Render("manage_article.tmpl")
}

type ArticlePublicController struct {
	tango.Ctx
	middle.AuthorizeRequire
	middle.Responsor
	middle.AdminRender
}

func (apc *ArticlePublicController) Get() {
	if id := apc.FormInt64("id"); id > 0 {
		if err := service.Call(service.Article.ToPublish, &id); err != nil {
			apc.RenderError(500, err)
			return
		}
	}
	apc.Redirect(apc.Req().Referer())
}

type ArticleDeleteController struct {
	tango.Ctx
	middle.AuthorizeRequire
	middle.AdminRender
}

func (adc *ArticleDeleteController) Get() {
	id := adc.FormInt64("id")
	if id > 0 {
		if err := service.Call(service.Article.Delete, id); err != nil {
			adc.RenderError(500, err)
			return
		}
	}
	adc.Redirect(adc.Req().Referer())
}
