package controller

import "pugo/src/middle"

type IndexController struct {
	middle.ThemeRender
}

func (ic *IndexController) Get() {
	ic.Render("index.tmpl")
}
