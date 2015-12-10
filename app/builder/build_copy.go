package builder

import (
	"os"
	"path"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo-static/app/helper"
)

// copy assets to target directory,
// favicon, robots.txt, error pages and all static asset if ctx.isCopyAllAssets
func (b *Builder) CopyAssets(ctx *Context) {
	// copy static files
	staticDir := ctx.Theme.Static()
	if err := b.copyAssets(ctx, staticDir, path.Join(ctx.DstDir, path.Base(staticDir))); err != nil {
		ctx.Error = err
		return
	}

	// copy media files
	if err := b.copyAssets(ctx, b.opt.MediaDir, path.Join(ctx.DstDir, path.Base(staticDir), path.Base(b.opt.MediaDir))); err != nil {
		ctx.Error = err
		return
	}

	// some extra files
	if b.copyExtraAssets(ctx); ctx.Error != nil {
		return
	}

	// clean un-track files
	if b.copyClean(ctx); ctx.Error != nil {
		return
	}
}

// clean old no change s file
func (b *Builder) copyClean(ctx *Context) {
	ctx.Error = filepath.Walk(ctx.DstDir, func(p string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		// not build file, clean it
		if !ctx.Diff.Exist(p) {
			ctx.Diff.Add(p, DIFF_REMOVE)
			return os.Remove(p)
		}
		return nil
	})
}

// copy static assets
func (b *Builder) copyAssets(ctx *Context, srcDir string, dstDir string) error {
	return filepath.Walk(srcDir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		// if file exist, check mod time
		rel, _ := filepath.Rel(srcDir, p)
		toFile := filepath.Join(dstDir, rel)
		if com.IsFile(toFile) {
			if fi2, _ := os.Stat(toFile); fi2 != nil {
				if fi2.ModTime().Sub(fi.ModTime()).Seconds() > 0 {
					ctx.Diff.Add(toFile, DIFF_KEEP)
					return nil
				}
				if err := helper.CopyFile(p, toFile); err != nil {
					return err
				}
				ctx.Diff.Add(toFile, DIFF_UPDATE)
				return nil
			}
		}

		// not exist, just copy
		if err := helper.CopyFile(p, toFile); err != nil {
			return err
		}
		ctx.Diff.Add(toFile, DIFF_ADD)
		return nil
	})
}

// copy extra files
func (b *Builder) copyExtraAssets(ctx *Context) {
	staticDir := ctx.Theme.Static()

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
		toFile := path.Join(ctx.DstOriginDir, f)
		if err := com.Copy(srcFile, toFile); err != nil {
			ctx.Error = err
			return
		}
		ctx.Diff.Add(toFile, DIFF_ADD)

		toFile = path.Join(ctx.DstDir, f)
		if err := com.Copy(srcFile, toFile); err != nil {
			ctx.Error = err
			return
		}
		ctx.Diff.Add(toFile, DIFF_ADD)
	}
}
