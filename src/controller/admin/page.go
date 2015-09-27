package admin

import (
	"github.com/fuxiaohei/pugo/src/middle"
	"github.com/fuxiaohei/pugo/src/model"
	"github.com/fuxiaohei/pugo/src/service"
	"github.com/fuxiaohei/pugo/src/utils"
	"github.com/lunny/tango"
	"github.com/tango-contrib/xsrf"
	"strings"
)

type PageWriteController struct {
	tango.Ctx
	xsrf.Checker

	middle.AuthorizeRequire
	middle.AdminRender
	middle.Validator
	middle.Responsor
}

func (pwc *PageWriteController) Get() {
	pwc.Title("WRITE PAGE - PUGO")
	pwc.Assign("XsrfHTML", pwc.XsrfFormHtml())
	if id := pwc.FormInt64("id"); id > 0 {
		var (
			opt  = service.PageReadOption{Id: id}
			page = new(model.Page)
		)
		if err := service.Call(service.Page.Read, opt, page); err != nil {
			pwc.RenderError(500, err)
			return
		}
		pwc.Assign("Page", page)
	}
	pwc.Render("write_page.tmpl")
}

// page post form
type PagePostForm struct {
	Title    string `form:"title" binding:"Required"`
	Link     string `form:"link" binding:"Required;AlphaDashDot"`
	Body     string `form:"body" binding:"Required"`
	Type     string `form:"type"`
	Draft    string `form:"draft"`
	Id       int64  `form:"id"`
	UserId   int64
	Comment  string `form:"comment"`
	Top      string `form:"top"`
	Template string `form:"template"`
}

func (f PagePostForm) toPage() *model.Page {
	page := &model.Page{
		Id:            f.Id,
		UserId:        f.UserId,
		Title:         f.Title,
		Link:          f.Link,
		Body:          f.Body,
		Status:        model.PAGE_STATUS_PUBLISH,
		CommentStatus: model.PAGE_COMMENT_OPEN,
		Hits:          1,
		Template:      f.Template,
	}
	switch strings.ToLower(f.Type) {
	case "html":
		page.BodyType = model.PAGE_BODY_HTML
	case "markdown":
		page.BodyType = model.PAGE_BODY_MARKDOWN
	default:
		page.BodyType = model.PAGE_BODY_MARKDOWN
	}
	if f.Draft == "draft" {
		page.Status = model.PAGE_STATUS_DRAFT
	}
	if f.Comment == "close" {
		page.CommentStatus = model.PAGE_COMMENT_CLOSE
	}
	if f.Top == "true" {
		page.TopLink = true
	}
	return page
}

func (pwc *PageWriteController) Post() {
	form := new(PagePostForm)
	if err := pwc.Validator.Validate(form, pwc); err != nil {
		pwc.JSONError(200, err)
		return
	}
	form.UserId = pwc.AuthUser.Id
	var page = new(model.Page)
	if err := service.Call(service.Page.Write, form.toPage(), page); err != nil {
		pwc.JSONError(200, err)
		return
	}
	pwc.JSON(map[string]interface{}{
		"page": page,
	})
}

type PageManageController struct {
	tango.Ctx
	middle.AuthorizeRequire
	middle.AdminRender
}

func (pmc *PageManageController) Get() {
	pmc.Title("PAGES - PUGO")
	var (
		opt = service.PageListOption{
			IsCount: true,
			Page:    pmc.FormInt("page", 0),
		}
		pages = make([]*model.Page, 0)
		pager = new(utils.Pager)
	)
	if err := service.Call(service.Page.List, opt, &pages, pager); err != nil {
		pmc.RenderError(500, err)
		return
	}
	pmc.Assign("Pages", pages)
	pmc.Assign("Pager", pager)
	pmc.Render("manage_page.tmpl")
}
