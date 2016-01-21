package builder

import (
	"os"
	"path"

	"github.com/go-xiaohei/pugo/app/model"
	"gopkg.in/inconshreveable/log15.v2"
)

// Compile compile source to static files
func Compile(ctx *Context) {
	var destDir = ""
	if destDir, ctx.Err = toDir(ctx.To); ctx.Err != nil {
		return
	}
	if ctx.Err = compilePosts(ctx, destDir); ctx.Err != nil {
		return
	}
}

func compilePosts(ctx *Context, toDir string) error {
	var (
		viewData map[string]interface{}
		dstFile  string
		err      error
	)
	for _, p := range ctx.Source.Posts {
		dstFile = path.Join(toDir, p.URL())
		if path.Ext(dstFile) == "" {
			dstFile += ".html"
		}

		viewData = ctx.View()
		viewData["Title"] = p.Title + " - " + ctx.Source.Meta.Title
		viewData["Desc"] = p.Desc
		viewData["Post"] = p
		viewData["Permalink"] = p.Permalink()
		viewData["PermaKey"] = p.Slug
		viewData["PostType"] = model.TreePost

		if err = compile(ctx, "post.html", viewData, dstFile); err != nil {
			return err
		}
	}
	return nil
}

func compile(ctx *Context, file string, viewData map[string]interface{}, destFile string) error {
	os.MkdirAll(path.Dir(destFile), os.ModePerm)
	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := ctx.Theme.Execute(f, file, viewData); err != nil {
		return err
	}
	log15.Debug("Build|To|%s", destFile)
	return nil
}
