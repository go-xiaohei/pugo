package builder

import (
	"path"
	"path/filepath"

	"github.com/go-xiaohei/pugo/app/model"
	"github.com/go-xiaohei/pugo/app/sync"
)

// Sync copy assets to destination directory
func Sync(ctx *Context) {
	if ctx.Err = ctx.Sync.SyncDir(ctx.Theme.StaticDir(), nil); ctx.Err != nil {
		return
	}

	opt := &sync.DirOption{
		Filter: func(p string) bool {
			return path.Ext(p) != ".md"
		},
	}
	var ignoreFiles []string

	opt.Prefix, _ = filepath.Rel(ctx.SrcDir(), ctx.SrcPostDir())
	files := model.ShouldPostMetaFiles()
	for _, f := range files {
		ignoreFiles = append(ignoreFiles, f)
	}
	opt.Ignore = ignoreFiles
	if ctx.Err = ctx.Sync.SyncDir(ctx.SrcPostDir(), opt); ctx.Err != nil {
		return
	}

	ignoreFiles = []string{}
	files = model.ShouldPageMetaFiles()
	for _, f := range files {
		ignoreFiles = append(ignoreFiles, f)
	}
	opt.Ignore = ignoreFiles
	opt.Prefix = ""
	if ctx.Err = ctx.Sync.SyncDir(ctx.SrcPageDir(), opt); ctx.Err != nil {
		return
	}

	opt.Prefix, _ = filepath.Rel(ctx.SrcDir(), ctx.SrcMediaDir())
	if ctx.Err = ctx.Sync.SyncDir(ctx.SrcMediaDir(), opt); ctx.Err != nil {
		return
	}
	opt.Ignore = []string{".git"}
	if ctx.Err = ctx.Sync.Clear(opt); ctx.Err != nil {
		return
	}
}
