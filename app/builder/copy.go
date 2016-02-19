package builder

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"strings"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/go-xiaohei/pugo/app/model"
	"gopkg.in/inconshreveable/log15.v2"
)

type (
	// CopiedOpt set option when copying files to destination
	CopiedOpt struct {
		MustFiles       []string
		IgnoreDir       []string
		IgnoreFile      []string
		CleanIgnoreDir  []string
		CleanIgnoreFile []string
	}
)

func defaultCopiedOpt() *CopiedOpt {
	return &CopiedOpt{
		MustFiles:       []string{"favicon.ico", "robots.txt"},
		IgnoreDir:       []string{".git"},
		IgnoreFile:      []string{"CNAME"},
		CleanIgnoreDir:  []string{".git"},
		CleanIgnoreFile: []string{"CNAME"},
	}
}

// Copy copy assets to destination directory
func Copy(ctx *Context) {
	if ctx.Err = CopyStatic(ctx); ctx.Err != nil {
		return
	}
	if ctx.Err = CopyMedia(ctx); ctx.Err != nil {
		return
	}
	if ctx.Err = CopyMust(ctx); ctx.Err != nil {
		return
	}
	if ctx.Err = CleanCopied(ctx); ctx.Err != nil {
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

		// ignore cases
		for _, ignoreDir := range ctx.Copied.IgnoreDir {
			if strings.HasPrefix(relPath, ignoreDir) {
				log15.Debug("Build|Ignore|%s", p)
				return nil
			}
		}
		for _, ignoreFile := range ctx.Copied.IgnoreFile {
			if relPath == ignoreFile {
				log15.Debug("Build|Ignore|%s", p)
				return nil
			}
		}

		if com.IsFile(toFile) {
			hash1, _ = helper.Md5File(p)
			hash2, _ = helper.Md5File(toFile)
			if hash1 == hash2 {
				ctx.Files.Add(toFile, fi.Size(), fi.ModTime(), model.FileStatic, model.OpKeep)
				log15.Debug("Build|Keep|%s", toFile)
				return nil
			}
		}

		// copy file
		os.MkdirAll(filepath.Dir(toFile), os.ModePerm)
		if err = com.Copy(p, toFile); err != nil {
			return err
		}

		ctx.Files.Add(toFile, fi.Size(), ctx.time, model.FileStatic, model.OpCopy)
		log15.Debug("Build|Copy|%s", toFile)

		return nil
	})
	return err
}

// CopyMust copy files in source to destination
func CopyMust(ctx *Context) error {
	var err error
	for _, f := range ctx.Copied.MustFiles {
		file := filepath.Join(ctx.SrcDir(), f)
		if com.IsFile(file) {
			toFile := filepath.Join(ctx.DstDir(), f)
			if err = com.Copy(file, toFile); err != nil {
				return err
			}
			ctx.Files.Add(toFile, 0, ctx.time, model.FileStatic, model.OpCopy)
			log15.Debug("Build|Copy|%s", toFile)
		}
	}
	return nil
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

// CleanCopied clean old copied files
func CleanCopied(ctx *Context) error {
	var (
		dstDir  = ctx.DstDir()
		relPath string
	)

	return filepath.Walk(dstDir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		relPath, _ = filepath.Rel(dstDir, p)

		for _, ignoreDir := range ctx.Copied.CleanIgnoreDir {
			if strings.HasPrefix(relPath, ignoreDir) {
				log15.Debug("Build|Ignore|%s", p)
				return nil
			}
		}
		for _, ignoreFile := range ctx.Copied.CleanIgnoreFile {
			if relPath == ignoreFile {
				log15.Debug("Build|Ignore|%s", p)
				return nil
			}
		}
		if !ctx.Files.Exist(p) {
			os.RemoveAll(p)
			ctx.Files.Add(p, fi.Size(), ctx.time, model.FileStatic, model.OpRemove)
			log15.Debug("Build|Remove|%s", p)
		}

		return nil
	})
}
