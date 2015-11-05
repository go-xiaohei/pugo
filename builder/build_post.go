package builder

import (
	"os"
	"path"
	"path/filepath"
	"pugo/model"
)

func (b *Builder) posts(ctx *context, r *Report) {
	if r.Error != nil {
		return
	}

	postDir := path.Join(b.srcDir, "post")
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
		post, err := model.NewPost(blocks, info)
		if err != nil {
			return err
		}
		if err = b.postRender(post, ctx); err != nil {
			return err
		}
		ctx.Posts = append(ctx.Posts, post)
		return f.Close()
	})
}

func (b *Builder) postRender(p *model.Post, ctx *context) error {
	template := b.Renders().Current().Template("post.html")
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

	viewData := map[string]interface{}{
		"Title": p.Title,
		"Post":  p,
	}
	if template.Compile(f, viewData, b.Renders().Current().FuncMap()); template.Error != nil {
		return err
	}
	return nil
}
