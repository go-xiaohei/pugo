package public

import (
	"github.com/go-xiaohei/pugo/src/middle"
	"github.com/go-xiaohei/pugo/src/model"
	"github.com/go-xiaohei/pugo/src/service"
	"github.com/lunny/tango"
)

type PageController struct {
	tango.Ctx
	middle.AuthorizeCheck
	middle.ThemeRender
}

func (pc *PageController) Get() {
	pageLink := pc.Param(":link")
	if len(pageLink) == 0 {
		pc.RenderError(404, nil)
		return
	}
	var (
		page = new(model.Page)
		opt  = service.PageReadOption{
			Id:        pc.ParamInt64(":id"),
			Link:      pageLink,
			Status:    model.PAGE_STATUS_PUBLISH,
			IsHit:     true,
			IsPublish: true,
		}
	)
	if err := service.Call(service.Page.Read, opt, page); err != nil {
		status := 500
		if err == service.ErrPageNotFound {
			status = 404
		}
		pc.RenderError(status, err)
		return
	}
	if page.Link != pageLink {
		pc.RenderError(404, nil)
		return
	}
	pc.Title(page.Title + " - " + service.Setting.General.Title)
	pc.Assign("Page", page)
	pc.Render(page.Template)
}
