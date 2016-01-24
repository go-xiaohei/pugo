package builder

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/go-xiaohei/pugo/app/model"
	"gopkg.in/inconshreveable/log15.v2"
)

// Copy copy assets to destination directory
func Copy(ctx *Context) {
	if ctx.Err = CopyStatic(ctx); ctx.Err != nil {
		return
	}
	if ctx.Err = CopyMedia(ctx); ctx.Err != nil {
		return
	}
}

func copyDirectory(ctx *Context, srcDir, dstDir string) error {
	var (
		toFile  string
		relPath string

		hash1, hash2 string
	)
	err := filepath.Walk(srcDir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		relPath, _ = filepath.Rel(srcDir, p)
		toFile = filepath.Join(dstDir, relPath)

		if com.IsFile(toFile) {
			hash1, _ = helper.Md5File(p)
			hash2, _ = helper.Md5File(toFile)
			if hash1 == hash2 {
				ctx.Files.Add(toFile, fi.Size(), fi.ModTime(), model.FileStatic)
				log15.Debug("Build|Keep|%s", toFile)
			}
			return nil
		}

		// copy file
		os.MkdirAll(filepath.Dir(toFile), os.ModePerm)
		if err = com.Copy(p, toFile); err != nil {
			return err
		}

		ctx.Files.Add(toFile, fi.Size(), ctx.time, model.FileStatic)
		log15.Debug("Build|Copy|%s", toFile)

		return nil
	})
	return err
}

// CopyStatic copy static assets from theme to destination directory
func CopyStatic(ctx *Context) error {
	if ctx.Theme == nil {
		return fmt.Errorf("CopyStatic need theme in Context")
	}
	return copyDirectory(ctx, ctx.Theme.StaticDir(), path.Join(ctx.dstDir, ctx.Source.Meta.Path))
}

// CopyMedia copy media files in source
func CopyMedia(ctx *Context) error {
	mediaDir := filepath.Join(ctx.srcDir, "media")
	if !com.IsDir(mediaDir) {
		return nil
	}
	return copyDirectory(ctx, mediaDir, filepath.Join(ctx.dstDir, ctx.Source.Meta.Path, "media"))
}
