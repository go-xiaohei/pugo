package server

import (
	"bytes"
	"github.com/Unknwon/com"
	"github.com/lunny/tango"
	"io/ioutil"
	"net/http"
	"path"
)

func Errors(dst string) tango.HandlerFunc {
	return func(ctx *tango.Context) {
		var (
			status     int
			errorBytes []byte
		)
		switch res := ctx.Result.(type) {
		case tango.AbortError:
			status = res.Code()
			errorBytes = []byte(res.Error())
		case error:
			status = http.StatusInternalServerError
			errorBytes = []byte(res.Error())
		default:
			status = http.StatusInternalServerError
			errorBytes = []byte(http.StatusText(http.StatusInternalServerError))
		}

		if status == 404 {
			notFoundFile := path.Join(dst, "errors/404.html")
			if com.IsFile(notFoundFile) {
				errorBytes, _ = ioutil.ReadFile(notFoundFile)
			}
			ctx.WriteHeader(status)
			ctx.ResponseWriter.Header().Add("Content-Type", "text/html")
			ctx.Write(errorBytes)
			return
		}

		if status == 500 {
			errorFile := path.Join(dst, "errors/500.html")
			if com.IsFile(errorFile) {
				fileBytes, _ := ioutil.ReadFile(errorFile)
				errorBytes = bytes.Replace(fileBytes, []byte("[error]"), errorBytes, -1)
			}
			ctx.WriteHeader(status)
			ctx.ResponseWriter.Header().Add("Content-Type", "text/html")
			ctx.Write(errorBytes)
			return
		}

		ctx.WriteHeader(status)
		ctx.Write(errorBytes)
	}
}
