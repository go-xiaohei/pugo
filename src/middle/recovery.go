package middle

import (
	"fmt"
	"github.com/lunny/tango"
	"gopkg.in/inconshreveable/log15.v2"
	"runtime"
)

func Recover() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		defer func() {
			if e := recover(); e != nil {
				header := fmt.Sprintf("Handler crashed with error: %v", e)
				content := header
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
			}
		}()

		ctx.Next()
	}
}
