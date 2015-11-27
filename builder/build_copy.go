package builder

import (
	"os"
	"path"
	"path/filepath"

	"strings"
	"time"

	"github.com/Unknwon/com"
)

// copy assets to target directory,
// favicon, robots.txt, error pages and all static asset if ctx.isCopyAllAssets
func (b *Builder) CopyAssets(ctx *Context) {
	if b.copyAssets(ctx); ctx.Error != nil {
		return
	}
	if b.copyError(ctx); ctx.Error != nil {
		return
	}
	if b.copyClean(ctx); ctx.Error != nil {
		return
	}
}

func (b *Builder) copyClean(ctx *Context) {
	filepath.Walk(ctx.DstDir, func(p string, info os.FileInfo, err error) error {
		if path.Ext(p) != ".html" ||
			info.IsDir() ||
			strings.HasPrefix(p, path.Join(ctx.DstDir, "/static")) ||
			strings.HasPrefix(p, path.Join(ctx.DstDir, "/upload")) {
			return nil
		}
		sub := info.ModTime().Sub(ctx.BeginTime.Add(-1 * time.Second))
		if sub < -1 {
			os.RemoveAll(p)
		}
		return nil
	})
}

// copy static assets
func (b *Builder) copyAssets(ctx *Context) {
	staticDir := ctx.Theme.Static()
	// copy all static
	dstDir := path.Join(ctx.DstDir, path.Base(staticDir))
	// remove old directory, otherwise return error when com.Copy
	os.RemoveAll(dstDir)
	if err := com.CopyDir(staticDir, dstDir); err != nil {
		ctx.Error = err
		return
	}

	assetFiles := []string{"favicon.ico", "robots.txt"}
	for _, f := range assetFiles {
		dstFile := path.Join(ctx.DstDir, f)
		srcFile := path.Join(staticDir, f)
		if err := com.Copy(srcFile, dstFile); err != nil {
			ctx.Error = err
			return
		}
	}
}

// copy error pages
func (b *Builder) copyError(ctx *Context) {
	// copy 404 to errors
	if err := b.copyErrorTemplate(ctx, "404.html"); err != nil {
		ctx.Error = err
		return
	}
	// copy 500 to errors
	if err := b.copyErrorTemplate(ctx, "500.html"); err != nil {
		ctx.Error = err
		return
	}
}

// copy error template,
// need render with some error data
func (b *Builder) copyErrorTemplate(ctx *Context, name string) error {
	dstFile := path.Join(ctx.DstDir, "error/"+name)
	os.MkdirAll(path.Dir(dstFile), os.ModePerm)
	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	return ctx.Theme.Execute(f, name, ctx.ViewData())
}
