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
	if err := removeDirectory(ctx.DstOriginDir); err != nil {
		ctx.Error = err
		return
	}
	if b.copyAssets(ctx); ctx.Error != nil {
		return
	}
}

// remove all sub dirs and files in directory
func removeDirectory(dir string) error {
	if !com.IsDir(dir) {
		return nil
	}
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
