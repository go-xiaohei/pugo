package builder

import (
	"github.com/Unknwon/com"
	"os"
	"path"
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
}

// copy static assets
func (b *Builder) copyAssets(ctx *Context, r *Report) {
	// copy all static
	if ctx.isCopyAllAssets {
		srcDir := b.Renders().Current().StaticDir()
		dstDir := path.Join(ctx.DstDir, path.Base(srcDir))
		com.CopyDir(srcDir, dstDir)
	}
	files := []string{"favicon.ico", "robots.txt"}
	staticDir := b.Renders().Current().StaticDir()
	for _, f := range files {
		srcFile := path.Join(staticDir, f)
		if com.IsFile(srcFile) {
			dstFile := path.Join(ctx.DstDir, f)
			com.Copy(srcFile, dstFile)
		}
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
	template := b.Renders().Current().Template(name)
	if template.Error != nil {
		return template.Error
	}
	dstFile := path.Join(ctx.DstDir, "errors/"+name)
	os.MkdirAll(path.Dir(dstFile), os.ModePerm)
	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	template.Compile(f, ctx.ViewData(), b.Renders().Current().FuncMap())
	return template.Error
}
