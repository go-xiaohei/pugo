package builder

import (
	"fmt"
	"github.com/go-xiaohei/pugo-static/model"
	"os"
	"path"
)

// compile data to html files
func (b *Builder) Compile(ctx *Context, r *Report) {
	if b.compileSinglePost(ctx, r); r.Error != nil {
		return
	}
	if b.compilePagedPost(ctx, r); r.Error != nil {
		return
	}
	if b.compileArchive(ctx, r); r.Error != nil {
		return
	}
	if b.compilePages(ctx, r); r.Error != nil {
		return
	}
	if b.compileTags(ctx, r); r.Error != nil {
		return
	}
	if b.compileIndex(ctx, r); r.Error != nil {
		return
	}
}

// compile each single post to html
func (b *Builder) compileSinglePost(ctx *Context, r *Report) {
	for _, p := range ctx.Posts {
		dstFile := path.Join(ctx.DstDir, p.Url)
		if path.Ext(dstFile) == "" {
			dstFile += ".html"
		}

		viewData := ctx.ViewData()
		viewData["Title"] = p.Title + " - " + ctx.Meta.Title
		viewData["Desc"] = p.Desc
		viewData["Post"] = p

		if err := b.compileTemplate("post.html", viewData, dstFile); err != nil {
			r.Error = err
			return
		}
	}
}

// compile paged posts to html
func (b *Builder) compilePagedPost(ctx *Context, r *Report) {
	// post pagination
	var (
		currentPosts []*model.Post = nil
		cursor                     = model.NewPagerCursor(4, len(ctx.Posts))
		page         int           = 1
		layout                     = "posts/%d"
	)
	for {
		pager := cursor.Page(page)
		if pager == nil {
			ctx.PostPageCount = page - 1
			break
		}

		currentPosts = ctx.Posts[pager.Begin:pager.End]
		pager.SetLayout("/" + layout)
		dstFile := path.Join(ctx.DstDir, fmt.Sprintf(layout+".html", pager.Page))

		viewData := ctx.ViewData()
		viewData["Title"] = fmt.Sprintf("Page %d - %s", pager.Page, ctx.Meta.Title)
		viewData["Posts"] = currentPosts
		viewData["Pager"] = pager

		if err := b.compileTemplate("posts.html", viewData, dstFile); err != nil {
			r.Error = err
			return
		}

		if page == 1 {
			ctx.indexPosts = currentPosts
			ctx.indexPager = pager
		}
		page++
	}
}

// compile archive page
func (b *Builder) compileArchive(ctx *Context, r *Report) {
	archives := model.NewArchive(ctx.Posts)
	dstFile := path.Join(ctx.DstDir, "archive.html")
	viewData := ctx.ViewData()
	viewData["Title"] = fmt.Sprintf("Archive - %s", ctx.Meta.Title)
	viewData["Archives"] = archives

	ctx.Navs.Hover("archive")
	defer ctx.Navs.Reset()

	if err := b.compileTemplate("archive.html", viewData, dstFile); err != nil {
		r.Error = err
		return
	}
}

// compile pages
func (b *Builder) compilePages(ctx *Context, r *Report) {
	for _, p := range ctx.Pages {
		dstFile := path.Join(ctx.DstDir, p.Url)
		if path.Ext(dstFile) == "" {
			dstFile += ".html"
		}

		ctx.Navs.Hover(p.HoverClass)
		defer ctx.Navs.Reset()
		viewData := ctx.ViewData()
		viewData["Title"] = p.Title + " - " + ctx.Meta.Title
		viewData["Desc"] = p.Desc
		viewData["Page"] = p

		if err := b.compileTemplate(p.Template, viewData, dstFile); err != nil {
			r.Error = err
			return
		}
	}
}

// compile tagged page
func (b *Builder) compileTags(ctx *Context, r *Report) {
	for t, posts := range ctx.tagPosts {
		dstFile := path.Join(ctx.DstDir, fmt.Sprintf("%s.html", ctx.Tags[t].Url))

		viewData := ctx.ViewData()
		viewData["Title"] = fmt.Sprintf("%s - %s", t, ctx.Meta.Title)
		viewData["Tag"] = ctx.Tags[t]
		viewData["Posts"] = posts

		if err := b.compileTemplate("posts.html", viewData, dstFile); err != nil {
			r.Error = err
			return
		}
	}
}

// compile index page
func (b *Builder) compileIndex(ctx *Context, r *Report) {
	template := "posts.html"
	if b.Renders().Current().IsExist("index.html") {
		template = "index.html"
	}

	dstFile := path.Join(ctx.DstDir, "index.html")

	ctx.Navs.Hover("home") // set hover status
	defer ctx.Navs.Reset() // remember to reset
	viewData := ctx.ViewData()
	viewData["Posts"] = ctx.indexPosts
	viewData["Pager"] = ctx.indexPager

	if err := b.compileTemplate(template, viewData, dstFile); err != nil {
		r.Error = err
		return
	}
}

// compile template by data and write to dest file.
func (b *Builder) compileTemplate(file string, viewData map[string]interface{}, destFile string) error {
	template := b.Renders().Current().Template(file)
	if template.Error != nil {
		return template.Error
	}
	os.MkdirAll(path.Dir(destFile), os.ModePerm)
	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	if template.Compile(f, viewData, b.Renders().Current().FuncMap()); template.Error != nil {
		return err
	}
	return nil
}
