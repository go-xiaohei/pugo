package public

import "github.com/go-xiaohei/pugo/src/middle"

type ErrorController struct {
	middle.ThemeRender
}

func (ec *ErrorController) NotFound() {
	ec.RenderError(404, nil)
}

func (ec *ErrorController) InternalError() {
	ec.RenderError(500, nil)
}
