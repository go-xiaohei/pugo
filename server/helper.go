package server

import (
	"github.com/Unknwon/com"
	"github.com/lunny/tango"
	"path"
	"pugo/builder"
	"strings"
)

type Helper struct {
	b      *builder.Builder
	dstDir string
}

func NewHelper(b *builder.Builder, dst string) *Helper {
	return &Helper{
		b:      b,
		dstDir: dst,
	}
}

func (h *Helper) Handle(ctx *tango.Context) {
	if h.b.IsBuilding() {
		ctx.Write([]byte("<h1>Pugo is building site!</h1>"))
		return
	}
	if h.b.Report().Error != nil {
		ctx.Abort(500, h.b.Report().Error.Error())
		return
	}
	// use built files
	url := ctx.Req().URL.Path
	// ignore static file
	if strings.HasPrefix(url, "/static") {
		ctx.Next()
		return
	}
	if url == "/" {
		url = "/index.html"
	}
	if path.Ext(url) == "" {
		url += ".html"
	}
	file := path.Join(h.dstDir, url)
	if com.IsFile(file) {
		ctx.ServeFile(file)
		return
	}

	ctx.Next()
}
