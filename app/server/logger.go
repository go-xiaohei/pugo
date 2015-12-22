package server

import (
	"fmt"
	"github.com/lunny/tango"
	"gopkg.in/inconshreveable/log15.v2"
	"path"
	"strings"
	"time"
)

var (
	logBase   = ""
	logFormat = "Http.%s.%d.%s" // Http.Method.Url from&duration
)

// log middleware handler
func logger() tango.HandlerFunc {
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
			if strings.HasPrefix(p, path.Join(logBase, "/static")) {
				return
			}
			if strings.HasPrefix(p, path.Join(logBase, "/theme")) {
				return
			}
			if p == "/favicon.ico" || p == "/robots.txt" {
				return
			}
		}

		statusCode := ctx.Status()
		if statusCode >= 200 && statusCode < 400 {
			log15.Info(
				fmt.Sprintf(logFormat, ctx.Req().Method, ctx.Status(), ctx.Req().URL.Path),
				"path", p,
				"remote", ctx.IP(),
				"duration", time.Since(start).Seconds()*1000,
			)
		} else if statusCode < 500 {
			log15.Warn(
				fmt.Sprintf(logFormat, ctx.Req().Method, ctx.Status(), ctx.Req().URL.Path),
				"path", p,
				"remote", ctx.IP(),
				"duration", time.Since(start).Seconds()*1000,
				"error", ctx.Result,
			)
		} else {
			log15.Error(
				fmt.Sprintf(logFormat, ctx.Req().Method, ctx.Status(), ctx.Req().URL.Path),
				"path", p,
				"remote", ctx.IP(),
				"duration", time.Since(start).Seconds()*1000,
				"error", ctx.Result,
			)
		}
	}
}
