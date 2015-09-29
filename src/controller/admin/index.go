package admin

import (
	"github.com/fuxiaohei/pugo/src/core"
	"github.com/fuxiaohei/pugo/src/middle"
	"github.com/fuxiaohei/pugo/src/model"
	"github.com/fuxiaohei/pugo/src/service"
	"strings"
)

type IndexController struct {
	middle.AuthorizeRequire
	middle.AdminRender
}

func (idx *IndexController) Get() {
	idx.Title(strings.ToUpper(core.PUGO_NAME))
	opt := service.MessageListOption{
		Page:    1,
		Size:    10,
		IsCount: false,
	}
	messages := make([]*model.Message, 0)
	if err := service.Call(service.Message.List, opt, &messages); err != nil {
		idx.RenderError(500, err)
		return
	}
	idx.Assign("Messages", messages)
	idx.Render("index.tmpl")
}
