package builder

import (
	"os"
	"path"
	"pugo/model"
)

func (b *Builder) nav(ctx *context, r *Report) {
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
