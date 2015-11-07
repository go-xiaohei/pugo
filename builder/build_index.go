package builder

import (
	"os"
	"path"
)

func (b *Builder) index(ctx *Context, r *Report) {
	if r.Error != nil {
		return
	}
	template := b.Renders().Current().Template("posts.html")
	if template.Error != nil {
		r.Error = template.Error
		return
	}
	dstFile := path.Join(ctx.DstDir, "index.html")
	os.MkdirAll(path.Dir(dstFile), os.ModePerm)
	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		r.Error = err
		return
	}
	defer f.Close()
	ctx.Navs.Hover("home") // set hover status
	defer ctx.Navs.Reset() // remember to reset
	viewData := ctx.ViewData()
	viewData["Posts"] = ctx.IndexPosts
	viewData["Pager"] = ctx.IndexPager
	if template.Compile(f, viewData, b.Renders().Current().FuncMap()); template.Error != nil {
		r.Error = err
		return
	}
}
