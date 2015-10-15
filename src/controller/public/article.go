package public

import (
	"fmt"
	"github.com/go-xiaohei/pugo/src/middle"
	"github.com/go-xiaohei/pugo/src/model"
	"github.com/go-xiaohei/pugo/src/service"
	"github.com/lunny/tango"
	"github.com/tango-contrib/session"
	"github.com/tango-contrib/xsrf"
)

type ArticleController struct {
	tango.Ctx
	xsrf.Checker
	session.Session

	middle.AuthorizeCheck
	middle.ThemeRender
}

func (ac *ArticleController) Get() {
	var (
		article = new(model.Article)
		opt     = service.ArticleReadOption{
			Id:        ac.ParamInt64(":id"),
			Link:      ac.Param(":link"),
			Status:    model.ARTICLE_STATUS_PUBLISH,
			IsHit:     true,
			IsPublish: true,
		}
		opt2 = service.CommentListOption{
			From:   model.COMMENT_FROM_ARTICLE,
			Status: 0,
		}
		comments = make([]*model.Comment, 0)
	)
	if err := service.Call(service.Article.Read, opt, article); err != nil {
		status := 500
		if err == service.ErrArticleNotFound {
			status = 404
		}
		ac.RenderError(status, err)
		return
	}
	if article.Id != opt.Id || article.Link != opt.Link {
		ac.RenderError(404, nil)
		return
	}
	opt2.FromId = article.Id
	if service.Setting.Comment.ShowWaitComment {
		opt2.AllAccessible = true
	} else {
		opt2.AllApproved = true
	}
	if err := service.Call(service.Comment.ListForContent, opt2, &comments); err != nil {
		ac.RenderError(500, err)
		return
	}

	shouldShowComments := true
	if article.IsCommentClosed() && len(comments) == 0 {
		shouldShowComments = false
	}

	if ac.AuthUser != nil {
		ac.Assign(middle.AuthUserTemplateField, nil) // set auth user nil instead of middleware assignment
		ac.Assign("FrontUser", model.NewFrontUser(ac.AuthUser))
	}

	ac.Title(article.Title + " - " + service.Setting.General.Title)
	ac.Assign("Article", article)
	ac.Assign("Comments", comments)
	ac.Assign("ShouldShowComments", shouldShowComments)
	ac.Assign("IsCommentEnable", article.IsCommentable(service.Setting.Comment.AutoCloseDay))
	ac.Assign("CommentUrl", fmt.Sprintf("/comment/article/%d", article.Id))
	ac.Assign("XsrfHTML", ac.XsrfFormHtml())
	ac.Render("single.tmpl")
}

type ArchiveController struct {
	tango.Ctx
	middle.ThemeRender
}

func (ac *ArchiveController) Get() {
	ac.Title("Archive - " + service.Setting.General.Title)
	var (
		opt      = service.ArchiveListOption{}
		archives = make([]*model.ArticleArchive, 0)
	)
	if err := service.Call(service.Article.Archive, opt, &archives); err != nil {
		ac.RenderError(500, err)
		return
	}
	ac.Assign("Archives", archives)
	ac.Render("archive.tmpl")
}
