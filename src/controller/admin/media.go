package admin

import (
	"pugo/src/middle"
)

type MediaController struct {
	middle.AuthorizeRequire
	middle.AdminRender
}

func (mc *MediaController) Get() {
	mc.Title("MEDIA - PUGO")
	mc.Render("manage_media.tmpl")
}
