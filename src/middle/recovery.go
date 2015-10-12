package middle

import (
	"fmt"
	"github.com/go-xiaohei/pugo/src/core"
	"github.com/lunny/tango"
	"gopkg.in/inconshreveable/log15.v2"
	"runtime"
)

func Recover() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		defer func() {
			if e := recover(); e != nil {
				header := fmt.Sprintf("%v", e)
				content := "Handler crashed with error:" + header
				for i := 1; ; i += 1 {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					} else {
						content += "\n"
					}
					content += fmt.Sprintf("%v %v", file, line)
				}

				p := ctx.Req().URL.Path
				if len(ctx.Req().URL.RawQuery) > 0 {
					p = p + "?" + ctx.Req().URL.RawQuery
				}

				if !ctx.Written() {
					ctx.Result = tango.InternalServerError(content)
					ctx.HandleError()
				}

				log15.Error(
					fmt.Sprintf(logFormat, ctx.Req().Method, ctx.Status(), ctx.Req().URL.Path),
					"path", p,
					"remote", ctx.Req().RemoteAddr,
					"error", header,
				)
				core.Crash.Error(
					fmt.Sprintf(logFormat, ctx.Req().Method, ctx.Status(), ctx.Req().URL.Path),
					"path", p,
					"remote", ctx.Req().RemoteAddr,
					"error", header,
				)
			}
		}()

		ctx.Next()
	}
}

type RecoveryHandler struct{}

func (rh *RecoveryHandler) Handle(ctx *tango.Context) {
	// capture render-controller error
	if render, ok := ctx.Action().(ITheme); ok {
		if err, ok := ctx.Result.(tango.AbortError); ok {
			render.RenderError(err.Code(), err)
			return
		}
		if err, ok := ctx.Result.(error); ok {
			ctx.WriteHeader(500)
			render.RenderError(ctx.Status(), err)
			return
		}
	}

	// capture abort error
	/*
		if err, ok := ctx.Result.(tango.AbortError); ok {
			ctx.WriteHeader(err.Code())
			theme := new(ThemeRender)
			theme.SetTheme(nil)
			theme.RenderError(err.Code(), err)
			return
		}*/

	// unexpected error
	tango.Errors()(ctx)
}
