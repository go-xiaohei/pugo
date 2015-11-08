package builder

import (
	"os"
	"path"
)

func (b *Builder) errors(ctx *Context, r *Report) {
	// copy 404 to errors
	if err := b.copyErrorTemplate(ctx, "404.html"); err != nil {
		r.Error = err
		return
	}
	// copy 500 to errors
	if err := b.copyErrorTemplate(ctx, "500.html"); err != nil {
		r.Error = err
		return
	}
}

func (b *Builder) copyErrorTemplate(ctx *Context, name string) error {
	template := b.Renders().Current().Template(name)
	if template.Error != nil {
		return template.Error
	}
	dstFile := path.Join(ctx.DstDir, "errors/"+name)
	os.MkdirAll(path.Dir(dstFile), os.ModePerm)
	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	template.Compile(f, ctx.ViewData(), b.Renders().Current().FuncMap())
	return template.Error
}
