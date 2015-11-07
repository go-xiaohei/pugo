package builder

import (
	"fmt"
	"os"
	"path"
	"pugo/model"
)

func (b *Builder) tags(ctx *context, r *Report) {
	if r.Error != nil {
		return
	}
	template := b.Renders().Current().Template("posts.html")
	if template.Error != nil {
		r.Error = template.Error
		return
	}
	// build tag pages
	var (
		tags     = make(map[string]*model.Tag)
		tagPosts = make(map[string][]*model.Post)
	)
	for _, p := range ctx.Posts {
		for i, t := range p.Tags {
			tags[t.Name] = &p.Tags[i]
			tagPosts[t.Name] = append(tagPosts[t.Name], p)
		}
	}
	ctx.Tags = tags
	for t, posts := range tagPosts {
		dstFile := path.Join(ctx.DstDir, fmt.Sprintf("%s.html", tags[t].Url))
		os.MkdirAll(path.Dir(dstFile), os.ModePerm)
		f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
		if err != nil {
			r.Error = err
			return
		}
		defer f.Close()

		viewData := ctx.viewData()
		viewData["Title"] = fmt.Sprintf("%s - %s", t, ctx.Meta.Title)
		viewData["Tag"] = tags[t]
		viewData["Posts"] = posts
		if template.Compile(f, viewData, b.Renders().Current().FuncMap()); template.Error != nil {
			r.Error = err
			return
		}
	}
}
