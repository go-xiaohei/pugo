package builder

import (
	"github.com/go-xiaohei/pugo-static/model"
	"os"
	"path"
	"path/filepath"
)

func (b *Builder) pages(ctx *Context, r *Report) {
	if r.Error != nil {
		return
	}

	// parse echo post
	postDir := path.Join(b.srcDir, "page")
	r.Error = filepath.Walk(postDir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if path.Ext(p) != ".md" {
			return nil
		}
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		blocks, err := b.parser.ParseReader(f)
		if err != nil {
			return err
		}
		page, err := model.NewPage(blocks, info)
		if err != nil {
			return err
		}
		if err = b.pageRender(page, ctx); err != nil {
			return err
		}
		ctx.Pages = append(ctx.Pages, page)
		return f.Close()
	})
	if r.Error != nil {
		return
	}
}

func (b *Builder) pageRender(p *model.Page, ctx *Context) error {
	template := b.Renders().Current().Template(p.Template)
	if template.Error != nil {
		return template.Error
	}
	dstFile := path.Join(ctx.DstDir, p.Url)
	if path.Ext(dstFile) == "" {
		dstFile += ".html"
	}

	os.MkdirAll(path.Dir(dstFile), os.ModePerm)

	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	ctx.Navs.Hover(p.HoverClass)
	defer ctx.Navs.Reset()
	viewData := ctx.ViewData()
	viewData["Title"] = p.Title + " - " + ctx.Meta.Title
	viewData["Desc"] = p.Desc
	viewData["Page"] = p
	if template.Compile(f, viewData, b.Renders().Current().FuncMap()); template.Error != nil {
		return err
	}
	return nil
}
