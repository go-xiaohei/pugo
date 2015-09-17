package public

import (
	"pugo/src/middle"
	"pugo/src/service"
)

type IndexController struct {
	middle.ThemeRender
}

func (ic *IndexController) Get() {
	ic.Title(service.Setting.General.FullTitle())
	ic.Render("index.tmpl")
}
