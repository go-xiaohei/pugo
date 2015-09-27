package admin

import (
	"github.com/fuxiaohei/pugo/src/core"
	"github.com/fuxiaohei/pugo/src/middle"
	"strings"
)

type IndexController struct {
	middle.AuthorizeRequire
	middle.AdminRender
}

func (idx *IndexController) Get() {
	idx.Title(strings.ToUpper(core.PUGO_NAME))
	idx.Render("index.tmpl")
}
