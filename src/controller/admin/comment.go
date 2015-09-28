package admin

import (
	"github.com/fuxiaohei/pugo/src/core"
	"github.com/fuxiaohei/pugo/src/middle"
	"github.com/fuxiaohei/pugo/src/model"
	"github.com/fuxiaohei/pugo/src/service"
	"github.com/fuxiaohei/pugo/src/utils"
	"github.com/lunny/tango"
	"strings"
)

type CommentController struct {
	tango.Ctx
	middle.AuthorizeRequire
	middle.AdminRender
}

func (cc *CommentController) Get() {
	cc.Title("COMMENT - " + strings.ToUpper(core.PUGO_NAME))
	var (
		opt = service.CommentListOption{
			Page:    cc.FormInt("page", 1),
			Size:    cc.FormInt("size", 10),
			IsCount: true,
		}
		comments = make([]*model.Comment, 0)
		pager    = new(utils.Pager)
	)

	// load comment
	switch cc.Form("status") {
	case "all":
		opt.Status = 0
	case "approved":
		opt.Status = model.COMMENT_STATUS_APPROVED
	case "wait":
		opt.Status = model.COMMENT_STATUS_WAIT
	case "spam":
		opt.Status = model.COMMENT_STATUS_SPAM
	default:
		opt.Status = 0
	}
	if err := service.Call(service.Comment.List, opt, &comments, pager); err != nil {
		cc.RenderError(500, err)
		return
	}

	// build pager url
	query := cc.Req().URL.Query()
	query.Del("page")
	queryStr := query.Encode()
	if len(queryStr) == 0 {
		queryStr = "/admin/manage/comment?page=%d"
	} else {
		queryStr = "/admin/manage/comment?" + queryStr + "&page=%d"
	}

	cc.Assign("PageUrl", queryStr)
	cc.Assign("Comments", comments)
	cc.Assign("Pager", pager)
	cc.Render("manage_comment.tmpl")
}

func (cc *CommentController) Approve() {
	if id := cc.FormInt64("id"); id > 0 {
		opt := service.CommentSwitchOption{
			Id:     id,
			Status: model.COMMENT_STATUS_APPROVED,
		}
		if err := service.Call(service.Comment.SwitchStatus, opt); err != nil {
			cc.RenderError(500, err)
			return
		}
	}
	cc.Redirect(cc.Req().Referer())
}

func (cc *CommentController) Delete() {
	if id := cc.FormInt64("id"); id > 0 {
		opt := service.CommentSwitchOption{
			Id:     id,
			Status: model.COMMENT_STATUS_DELETED,
		}
		if err := service.Call(service.Comment.SwitchStatus, opt); err != nil {
			cc.RenderError(500, err)
			return
		}
	}
	cc.Redirect(cc.Req().Referer())
}

func (cc *CommentController) Reply() {
	c := &model.Comment{
		UserId:   cc.AuthUser.Id,
		Body:     cc.Form("content"),
		ParentId: cc.FormInt64("pid"),
	}
	if err := service.Call(service.Comment.Reply, c); err != nil {
		cc.RenderError(500, err)
		return
	}
	cc.Redirect("/admin/manage/comment")
}
