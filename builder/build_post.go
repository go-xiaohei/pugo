package builder

import (
	"fmt"
	"github.com/go-xiaohei/pugo-static/model"
	"os"
	"path"
	"path/filepath"
	"sort"
)

func (b *Builder) posts(ctx *Context, r *Report) {
	if r.Error != nil {
		return
	}

	// parse echo post
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
	if r.Error != nil {
		return
	}

	// parse all post
	sort.Sort(model.Posts(ctx.Posts))
	var (
		currentPosts []*model.Post = nil
		cursor                     = model.NewPagerCursor(4, len(ctx.Posts))
		page         int           = 1
	)
	for {
		pager := cursor.Page(page)
		if pager == nil {
			ctx.PostPages = page - 1
			break
		}
		currentPosts = ctx.Posts[pager.Begin:pager.End]
		if err := b.postsRender(currentPosts, ctx, pager); err != nil {
			r.Error = err
			return
		}

		if page == 1 {
			ctx.IndexPosts = currentPosts
			ctx.IndexPager = pager
		}
		page++
	}

	// build archive
	archives := model.NewArchive(ctx.Posts)
	if err := b.archiveRender(archives, ctx); err != nil {
		r.Error = err
		return
	}
}

func (b *Builder) postRender(p *model.Post, ctx *Context) error {
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

	viewData := ctx.ViewData()
	viewData["Title"] = p.Title + " - " + ctx.Meta.Title
	viewData["Desc"] = p.Desc
	viewData["Post"] = p
	if template.Compile(f, viewData, b.Renders().Current().FuncMap()); template.Error != nil {
		return err
	}
	return nil
}

func (b *Builder) postsRender(posts []*model.Post, ctx *Context, pager *model.Pager) error {
	template := b.Renders().Current().Template("posts.html")
	if template.Error != nil {
		return template.Error
	}
	layout := "posts/%d"
	pager.SetLayout("/" + layout)

	dstFile := path.Join(ctx.DstDir, fmt.Sprintf(layout+".html", pager.Page))
	os.MkdirAll(path.Dir(dstFile), os.ModePerm)
	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	viewData := ctx.ViewData()
	viewData["Title"] = fmt.Sprintf("Page %d - %s", pager.Page, ctx.Meta.Title)
	viewData["Posts"] = posts
	viewData["Pager"] = pager
	if template.Compile(f, viewData, b.Renders().Current().FuncMap()); template.Error != nil {
		return err
	}
	return nil
}

func (b *Builder) archiveRender(archives []*model.Archive, ctx *Context) error {
	template := b.Renders().Current().Template("archive.html")
	if template.Error != nil {
		return template.Error
	}

	dstFile := path.Join(ctx.DstDir, "archive.html")
	os.MkdirAll(path.Dir(dstFile), os.ModePerm)
	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	ctx.Navs.Hover("archive")
	defer ctx.Navs.Reset()

	viewData := ctx.ViewData()
	viewData["Title"] = fmt.Sprintf("Archive - %s", ctx.Meta.Title)
	viewData["Archives"] = archives
	if template.Compile(f, viewData, b.Renders().Current().FuncMap()); template.Error != nil {
		return err
	}
	return nil
}
