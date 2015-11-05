package builder

import (
	"os"
	"path"
)

func (b *Builder) index(ctx *context, r *Report) {
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

	viewData := map[string]interface{}{
		"Nav":   ctx.Navs,
		"Posts": ctx.IndexPosts,
		"Pager": ctx.IndexPager,
	}
	if template.Compile(f, viewData, b.Renders().Current().FuncMap()); template.Error != nil {
		r.Error = err
		return
	}
}
