package middle

import (
	"fmt"
	"github.com/lunny/tango"
	"gopkg.in/inconshreveable/log15.v2"
	"strings"
	"time"
)

var (
	logFormat = "Http.%s.%d.%s" // Http.Method.Url from&duration
)

// log middleware handler
func Logger() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		start := time.Now()
		p := ctx.Req().URL.Path
		if len(ctx.Req().URL.RawQuery) > 0 {
			p = p + "?" + ctx.Req().URL.RawQuery
		}
		if action := ctx.Action(); action != nil {
			if l, ok := action.(tango.LogInterface); ok {
				l.SetLogger(ctx.Logger)
			}
		}
		ctx.Next()
		if !ctx.Written() {
			if ctx.Result == nil {
				ctx.Result = tango.NotFound()
			}
			ctx.HandleError()
		}

		// skip static files
		if ctx.Req().Method == "GET" {
			if strings.HasPrefix(p, "/static") {
				return
			}
		}

		statusCode := ctx.Status()
		if statusCode >= 200 && statusCode < 400 {
			log15.Info(
				fmt.Sprintf(logFormat, ctx.Req().Method, ctx.Status(), ctx.Req().URL.Path),
				"path", p,
				"remote", ctx.Req().RemoteAddr,
				"duration", time.Since(start).Seconds()*1000,
			)
		} else {
			log15.Error(
				fmt.Sprintf(logFormat, ctx.Req().Method, ctx.Status(), ctx.Req().URL.Path),
				"path", p,
				"remote", ctx.Req().RemoteAddr,
				"duration", time.Since(start).Seconds()*1000,
				"error", ctx.Result,
			)
		}
	}
}

// return friendly remote string
func friendRemoteString(remote string) string {
	if len(remote) < 21 {
		remote += strings.Repeat(" ", 21-len(remote))
	}
	return remote
}
