package admin

import (
	"github.com/lunny/tango"
	"pugo/src/core"
	"pugo/src/middle"
	"pugo/src/model"
	"pugo/src/service"
	"pugo/src/utils"
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
	cc.Assign("Comments", comments)
	cc.Assign("Pager", pager)
	cc.Render("manage_comment.tmpl")
}
