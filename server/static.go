package server

import (
	"github.com/lunny/tango"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Static struct {
	RootPath   string
	Prefix     string
	IndexFiles []string
	ListDir    bool
	FilterExts []string
}

func NewStatic() *Static {
	return &Static{
		RootPath:   "./static",
		Prefix:     "/static",
		IndexFiles: []string{"index.html", "index.htm"},
		ListDir:    false,
		FilterExts: nil,
	}
}

func (s *Static) Handle(ctx *tango.Context) {
	if ctx.Req().Method != "GET" && ctx.Req().Method != "HEAD" {
		ctx.Next()
		return
	}
	var rPath = ctx.Req().URL.Path
	// if defined prefix, then only check prefix
	if s.Prefix != "" {
		if !strings.HasPrefix(ctx.Req().URL.Path, s.Prefix) {
			ctx.Next()
			return
		} else {
			if len(s.Prefix) == len(ctx.Req().URL.Path) {
				rPath = ""
			} else {
				rPath = ctx.Req().URL.Path[len(s.Prefix):]
			}
		}
	}

	fPath, _ := filepath.Abs(filepath.Join(s.RootPath, rPath))
	finfo, err := os.Stat(fPath)
	if err != nil {
		if !os.IsNotExist(err) {
			ctx.Result = tango.InternalServerError(err.Error())
			ctx.HandleError()
			return
		}
	} else if !finfo.IsDir() {
		if len(s.FilterExts) > 0 {
			var matched bool
			for _, ext := range s.FilterExts {
				if filepath.Ext(fPath) == ext {
					matched = true
					break
				}
			}
			if !matched {
				ctx.Next()
				return
			}
		}

		err := ctx.ServeFile(fPath)
		if err != nil {
			ctx.Result = tango.InternalServerError(err.Error())
			ctx.HandleError()
		}
		return
	} else {
		// try serving index.html or index.htm
		if len(s.IndexFiles) > 0 {
			for _, index := range s.IndexFiles {
				nPath := filepath.Join(fPath, index)
				finfo, err = os.Stat(nPath)
				if err != nil {
					if !os.IsNotExist(err) {
						ctx.Result = tango.InternalServerError(err.Error())
						ctx.HandleError()
						return
					}
				} else if !finfo.IsDir() {
					err = ctx.ServeFile(nPath)
					if err != nil {
						ctx.Result = tango.InternalServerError(err.Error())
						ctx.HandleError()
					}
					return
				}
			}
		}

		// list dir files
		if s.ListDir {
			ctx.Header().Set("Content-Type", "text/html; charset=UTF-8")
			ctx.Write([]byte(`<ul style="list-style-type:none;line-height:32px;">`))
			rootPath, _ := filepath.Abs(s.RootPath)
			rPath, _ := filepath.Rel(rootPath, fPath)
			if fPath != rootPath {
				ctx.Write([]byte(`<li>&nbsp; &nbsp; <a href="/` + path.Join(s.Prefix, filepath.Dir(rPath)) + `">..</a></li>`))
			}
			err = filepath.Walk(fPath, func(p string, fi os.FileInfo, err error) error {
				rPath, _ := filepath.Rel(fPath, p)
				if rPath == "." || len(strings.Split(rPath, string(filepath.Separator))) > 1 {
					return nil
				}
				rPath, _ = filepath.Rel(rootPath, p)
				ps, _ := os.Stat(p)
				if ps.IsDir() {
					ctx.Write([]byte(`<li>┖ <a href="/` + path.Join(s.Prefix, rPath) + `">` + filepath.Base(p) + `</a></li>`))
				} else {
					if len(s.FilterExts) > 0 {
						var matched bool
						for _, ext := range s.FilterExts {
							if filepath.Ext(p) == ext {
								matched = true
								break
							}
						}
						if !matched {
							return nil
						}
					}

					ctx.Write([]byte(`<li>&nbsp; &nbsp; <a href="/` + path.Join(s.Prefix, rPath) + `">` + filepath.Base(p) + `</a></li>`))
				}
				return nil
			})
			ctx.Write([]byte("</ul>"))
			return
		}
	}

	ctx.Next()
}
