package server

import (
	"github.com/lunny/tango"
	"pugo/builder"
)

type Helper struct {
	b *builder.Builder
}

func NewHelper(b *builder.Builder) *Helper {
	return &Helper{b: b}
}

func (h *Helper) Handle(ctx *tango.Context) {
	if h.b.IsBuilding() {
		ctx.Write([]byte("<h1>Pugo is building site!</h1>"))
		return
	}
	if h.b.Report().Error != nil {
		ctx.WriteHeader(500)
		ctx.Write([]byte("<h1>Pugo built fail:" + h.b.Report().Error.Error() + "</h1>"))
		return
	}
	ctx.Next()
}
