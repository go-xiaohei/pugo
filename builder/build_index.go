package builder

import (
	"os"
	"path"
)

func (b *Builder) index(r *Report) {
	if r.Error != nil {
		return
	}
	template := b.Renders().Current().Template("index.html")
	if template.Error != nil {
		r.Error = template.Error
		return
	}
	dstFile := path.Join(r.DstDir, "index.html")
	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		r.Error = err
		return
	}
	defer f.Close()
	if template.Compile(f); template.Error != nil {
		r.Error = template.Error
		return
	}
}
