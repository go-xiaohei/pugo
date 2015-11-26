package builder

import (
	"os"
	"path"
	"path/filepath"

	"github.com/Unknwon/com"
)

// copy assets to target directory,
// favicon, robots.txt, error pages and all static asset if ctx.isCopyAllAssets
func (b *Builder) CopyAssets(ctx *Context, r *Report) {
	if b.copyAssets(ctx, r); r.Error != nil {
		return
	}
	if b.copyError(ctx, r); r.Error != nil {
		return
	}
	if b.copyClean(ctx, r); r.Error != nil {
		return
	}
}

func (b *Builder) copyClean(ctx *Context, r *Report) {
	filepath.Walk(ctx.DstDir, func(p string, info os.FileInfo, err error) error {
		return nil
	})
}

// copy static assets
func (b *Builder) copyAssets(ctx *Context, r *Report) {
	staticDir := ctx.Theme.Static()
	// copy all static
	dstDir := path.Join(ctx.DstDir, path.Base(staticDir))
	// remove old directory, otherwise return error when com.Copy
	os.RemoveAll(dstDir)
	if err := com.CopyDir(staticDir, dstDir); err != nil {
		r.Error = err
		return
	}
}

// copy error pages
func (b *Builder) copyError(ctx *Context, r *Report) {
	// copy 404 to errors
	if err := b.copyErrorTemplate(ctx, "404.html"); err != nil {
		r.Error = err
		return
	}
	// copy 500 to errors
	if err := b.copyErrorTemplate(ctx, "500.html"); err != nil {
		r.Error = err
		return
	}
}

// copy error template,
// need render with some error data
func (b *Builder) copyErrorTemplate(ctx *Context, name string) error {
	dstFile := path.Join(ctx.DstDir, "errors/"+name)
	os.MkdirAll(path.Dir(dstFile), os.ModePerm)
	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	return ctx.Theme.Execute(f, name, ctx.ViewData())
}
