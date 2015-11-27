package builder

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/Unknwon/com"
)

// copy assets to target directory,
// favicon, robots.txt, error pages and all static asset if ctx.isCopyAllAssets
func (b *Builder) CopyAssets(ctx *Context) {
	if b.copyClean(ctx); ctx.Error != nil {
		return
	}
	if b.copyAssets(ctx); ctx.Error != nil {
		return
	}
	if b.copyError(ctx); ctx.Error != nil {
		return
	}
}

func (b *Builder) copyClean(ctx *Context) {
	if err := removeDirectory(ctx.DstOriginDir); err != nil {
		ctx.Error = err
	}
}

func removeDirectory(dir string) error {
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, d := range dirs {
		if d.IsDir() {
			if err = removeDirectory(path.Join(dir, d.Name())); err != nil {
				return err
			}
		}
	}
	return os.RemoveAll(dir)
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
		// use origin dir, make these files existing in top directory
		dstFile := path.Join(ctx.DstOriginDir, f)
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
