package builder

import (
	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo-static/model"
	"os"
	"path"
)

func (b *Builder) meta(ctx *Context, r *Report) {
	if r.Error != nil {
		return
	}
	navFile := path.Join(b.srcDir, "meta.md")
	f, err := os.Open(navFile)
	if err != nil {
		r.Error = err
		return
	}
	blocks, err := b.parser.ParseReader(f)
	if err != nil {
		r.Error = err
		return
	}
	meta, err := model.NewMeta(blocks)
	if err != nil {
		r.Error = err
		return
	}
	ctx.Meta = meta
}

func (b *Builder) nav(ctx *Context, r *Report) {
	if r.Error != nil {
		return
	}
	navFile := path.Join(b.srcDir, "nav.md")
	f, err := os.Open(navFile)
	if err != nil {
		r.Error = err
		return
	}
	blocks, err := b.parser.ParseReader(f)
	if err != nil {
		r.Error = err
		return
	}
	navs, err := model.NewNavs(blocks)
	if err != nil {
		r.Error = err
		return
	}
	ctx.Navs = navs
}

func (b *Builder) comment(ctx *Context, r *Report) {
	if r.Error != nil {
		return
	}
	cmtFile := path.Join(b.srcDir, "comment.md")
	f, err := os.Open(cmtFile)
	if err != nil {
		r.Error = err
		return
	}
	blocks, err := b.parser.ParseReader(f)
	if err != nil {
		r.Error = err
		return
	}
	cmt, err := model.NewComment(blocks)
	if err != nil {
		r.Error = err
		return
	}
	ctx.Comment = cmt
}

func (b *Builder) assets(ctx *Context, r *Report) {
	files := []string{"favicon.ico", "robot.txt"}
	staticDir := b.Renders().Current().StaticDir()
	for _, f := range files {
		srcFile := path.Join(staticDir, f)
		if com.IsFile(srcFile) {
			dstFile := path.Join(ctx.DstDir, f)
			com.Copy(srcFile, dstFile)
		}
	}
}
