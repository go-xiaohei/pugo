package builder

import (
	"path"
	"path/filepath"

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
	opt.Prefix, _ = filepath.Rel(ctx.SrcDir(), ctx.SrcPostDir())
	if ctx.Err = ctx.Sync.SyncDir(ctx.SrcPostDir(), opt); ctx.Err != nil {
		return
	}
	opt.Prefix, _ = filepath.Rel(ctx.SrcDir(), ctx.SrcPageDir())
	if ctx.Err = ctx.Sync.SyncDir(ctx.SrcPageDir(), opt); ctx.Err != nil {
		return
	}
	opt.Prefix, _ = filepath.Rel(ctx.SrcDir(), ctx.SrcMediaDir())
	if ctx.Err = ctx.Sync.SyncDir(ctx.SrcMediaDir(), opt); ctx.Err != nil {
		return
	}
	if ctx.Err = ctx.Sync.Clear(); ctx.Err != nil {
		return
	}
}
