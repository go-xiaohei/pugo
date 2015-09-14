package admin

import (
	"github.com/tango-contrib/xsrf"
	"pugo/src/middle"
	"pugo/src/model"
	"pugo/src/service"
	"strings"
)

type ArticleWriteController struct {
	xsrf.Checker

	middle.AuthorizeRequire
	middle.AdminRender
	middle.Validator
	middle.Responsor
}

func (awc *ArticleWriteController) Get() {
	awc.Title("WRITE ARTICLE - PUGO")
	awc.Assign("XsrfHTML", awc.XsrfFormHtml())
	awc.Render("write_article.tmpl")
}

// article post form
type ArticleForm struct {
	Title  string `form:"title" binding:"Required"`
	Link   string `form:"link" binding:"Required;AlphaDashDot"`
	Body   string `form:"body" binding:"Required"`
	Type   string `form:"type"`
	Tag    string `form:"tag"`
	Draft  string `form:"draft"`
	Id     int64  `form:"id"`
	UserId int64
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
	}
	if strings.Contains(f.Body, "<!--more-->") {
		article.Preview = strings.Split(f.Body, "<!--more-->")[0]
	}
	switch strings.ToLower(f.Type) {
	case "html":
		article.BodyType = model.ARTICLE_BODY_HTML
	case "markdown":
	default:
		article.BodyType = model.ARTICLE_BODY_MARKDOWN
	}
	if f.Draft == "draft" {
		article.Status = model.ARTICLE_STATUS_DRAFT
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
	middle.AuthorizeRequire
	middle.AdminRender
}

func (amc *ArticleManageController) Get() {
	amc.Title("ARTICLES - PUGO")
	amc.Render("manage_article.tmpl")
}
