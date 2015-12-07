package builder

import (
	"os"
	"path"

	"path/filepath"
	"strings"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo-static/app/helper"
)

// copy assets to target directory,
// favicon, robots.txt, error pages and all static asset if ctx.isCopyAllAssets
func (b *Builder) CopyAssets(ctx *Context) {
	if b.copyAssets(ctx); ctx.Error != nil {
		return
	}
	if b.copyClean(ctx); ctx.Error != nil {
		return
	}

}

// clean old no change s file
func (b *Builder) copyClean(ctx *Context) {
	static := path.Base(ctx.Theme.Static())
	ctx.Error = filepath.Walk(ctx.DstDir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// skip directory and static dir
		rel, _ := filepath.Rel(ctx.DstDir, p)
		if fi.IsDir() || strings.HasPrefix(rel, static) {
			return nil
		}
		if rel == "favicon.ico" || rel == "robots.txt" {
			return nil
		}
		ext := path.Ext(p)
		if ext == ".html" || ext == ".xml" {
			if ctx.BeginTime.Unix()-fi.ModTime().Unix() > 10 {
				return os.Remove(p)
			}
		}
		return nil
	})
}

// copy static assets
func (b *Builder) copyAssets(ctx *Context) {
	staticDir := ctx.Theme.Static()
	// copy all static
	dstDir := path.Join(ctx.DstDir, path.Base(staticDir))
	if err := helper.RemoveDir(dstDir); err != nil {
		ctx.Error = err
		return
	}
	if err := com.CopyDir(staticDir, dstDir, func(p string) bool {
		if path.Ext(p) == ".DS_Store" {
			return true
		}
		return false
	}); err != nil {
		ctx.Error = err
		return
	}

	// copy upload data
	if com.IsDir(b.opt.MediaDir) {
		dstDir = path.Join(ctx.DstDir, path.Base(staticDir), path.Base(b.opt.MediaDir))
		if !com.IsDir(dstDir) {
			if err := com.CopyDir(b.opt.MediaDir, dstDir); err != nil {
				ctx.Error = err
				return
			}
		}
	}

	assetFiles := []string{"favicon.ico", "robots.txt"}
	for _, f := range assetFiles {
		srcFile := path.Join(b.opt.SrcDir, f)
		if !com.IsFile(srcFile) {
			srcFile = path.Join(staticDir, f)
		}
		if !com.IsFile(srcFile) {
			continue
		}

		// use origin dir, make these files existing in top directory
		if err := com.Copy(srcFile, path.Join(ctx.DstOriginDir, f)); err != nil {
			ctx.Error = err
			return
		}
		if err := com.Copy(srcFile, path.Join(ctx.DstDir, f)); err != nil {
			ctx.Error = err
			return
		}
	}
}
