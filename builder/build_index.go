package builder

import (
	"os"
	"path"
)

func (b *Builder) index(ctx *context, r *Report) {
	if r.Error != nil {
		return
	}
	template := b.Renders().Current().Template("index.html")
	if template.Error != nil {
		r.Error = template.Error
		return
	}
	dstFile := path.Join(ctx.DstDir, "index.html")
	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		r.Error = err
		return
	}
	defer f.Close()
	if template.Compile(f, nil, b.Renders().Current().FuncMap()); template.Error != nil {
		r.Error = template.Error
		return
	}
}
