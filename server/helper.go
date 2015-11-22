package server

import (
	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo-static/builder"
	"github.com/lunny/tango"
	"gopkg.in/inconshreveable/log15.v2"
	"path"
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
	if path.Ext(url) == "" && url[len(url)-1:] != "/" {
		url += "/"
	}

	// fix url if base is not empty
	base := h.b.Context().Meta.Base

	if base != "" {
		if !strings.HasPrefix(url, base) {
			ctx.NotFound()
			return
		}
		url = strings.TrimPrefix(url, base)
	}

	// ignore static file
	if strings.HasPrefix(url, "/static") {
		ctx.Next()
		return
	}
	// ignore upload file
	if strings.HasPrefix(url, "/static") {
		ctx.Next()
		return
	}

	if url == "/" {
		url = "/index.html"
	}
	if path.Ext(url) == "" {
		url = strings.TrimSuffix(url, "/")
		url += ".html"
	}
	file := path.Join(h.dstDir, url)
	log15.Debug("Server.StaticFile."+file, "request", ctx.Req().URL.Path)
	if com.IsFile(file) {
		ctx.ServeFile(file)
		return
	}

	ctx.Next()
}
