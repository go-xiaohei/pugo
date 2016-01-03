package builder

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/go-xiaohei/pugo/app/model"
	"gopkg.in/inconshreveable/log15.v2"
)

// Compile builds data to html files
func (b *Builder) Compile(ctx *Context) {
	if b.compileSinglePost(ctx); ctx.Error != nil {
		return
	}
	if b.compilePagedPost(ctx); ctx.Error != nil {
		return
	}
	if b.compileArchive(ctx); ctx.Error != nil {
		return
	}
	if b.compilePages(ctx); ctx.Error != nil {
		return
	}
	if b.compileTags(ctx); ctx.Error != nil {
		return
	}
	if b.compileIndex(ctx); ctx.Error != nil {
		return
	}
}

// compile each single post to html
func (b *Builder) compileSinglePost(ctx *Context) {
	log15.Debug("Post." + strconv.Itoa(len(ctx.Posts)))
	for _, p := range ctx.Posts {
		dstFile := path.Join(ctx.DstDir, p.Url)
		if path.Ext(dstFile) == "" {
			dstFile += ".html"
		}

		viewData := ctx.ViewData()
		viewData["Title"] = p.Title + " - " + ctx.Meta.Title
		viewData["Desc"] = p.Desc
		viewData["Post"] = p
		viewData["Permalink"] = p.Permalink
		viewData["PostType"] = "post"

		if err := b.compileTemplate(ctx, "post.html", viewData, dstFile); err != nil {
			ctx.Error = err
			return
		}
	}
}

// compile paged posts to html
func (b *Builder) compilePagedPost(ctx *Context) {
	// post pagination
	var (
		currentPosts []*model.Post
		cursor       = helper.NewPagerCursor(4, len(ctx.Posts))
		page         = 1
		layout       = "posts/%d"
	)
	for {
		pager := cursor.Page(page)
		if pager == nil {
			ctx.PostPageCount = page - 1
			break
		}

		currentPosts = ctx.Posts[pager.Begin:pager.End]
		pager.SetLayout("/" + layout + ".html")

		dstFile := path.Join(ctx.DstDir, fmt.Sprintf(layout+".html", pager.Current))

		viewData := ctx.ViewData()
		viewData["Title"] = fmt.Sprintf("Page %d - %s", pager.Current, ctx.Meta.Title)
		viewData["Posts"] = currentPosts
		viewData["Pager"] = pager
		viewData["PostType"] = "post-list"

		if err := b.compileTemplate(ctx, "posts.html", viewData, dstFile); err != nil {
			ctx.Error = err
			return
		}

		if page == 1 {
			ctx.indexPosts = currentPosts
			ctx.indexPager = pager
		}
		page++
	}
	log15.Debug("Post.Pages." + strconv.Itoa(page-1))
}

// compile archive page
func (b *Builder) compileArchive(ctx *Context) {
	archives := model.NewArchive(ctx.Posts)
	dstFile := path.Join(ctx.DstDir, "archive.html")
	viewData := ctx.ViewData()
	viewData["Title"] = fmt.Sprintf("Archive - %s", ctx.Meta.Title)
	viewData["Archives"] = archives
	viewData["PostType"] = "post-archive"

	ctx.Navs.Hover("archive")
	defer ctx.Navs.Reset()

	if err := b.compileTemplate(ctx, "archive.html", viewData, dstFile); err != nil {
		ctx.Error = err
		return
	}
}

// compile pages
func (b *Builder) compilePages(ctx *Context) {
	log15.Debug("Pages." + strconv.Itoa(len(ctx.Pages)))
	for _, p := range ctx.Pages {
		dstFile := path.Join(ctx.DstDir, p.Url)
		if path.Ext(dstFile) == "" {
			dstFile += ".html"
		}

		ctx.Navs.Hover(p.HoverClass)
		viewData := ctx.ViewData()
		viewData["Title"] = p.Title + " - " + ctx.Meta.Title
		viewData["Desc"] = p.Desc
		viewData["Page"] = p
		viewData["Permalink"] = p.Permalink
		viewData["PostType"] = "page"

		if err := b.compileTemplate(ctx, p.Template, viewData, dstFile); err != nil {
			ctx.Error = err
			ctx.Navs.Reset()
			return
		}
		ctx.Navs.Reset()
	}
}

// compile tagged page
func (b *Builder) compileTags(ctx *Context) {
	log15.Debug("Tags." + strconv.Itoa(len(ctx.Tags)))
	for t, posts := range ctx.tagPosts {
		dstFile := path.Join(ctx.DstDir, fmt.Sprintf("tags/%s.html", ctx.Tags[t].Name))

		viewData := ctx.ViewData()
		viewData["Title"] = fmt.Sprintf("%s - %s", t, ctx.Meta.Title)
		viewData["Tag"] = ctx.Tags[t]
		viewData["Posts"] = posts
		viewData["PostType"] = "post-tag"

		if err := b.compileTemplate(ctx, "posts.html", viewData, dstFile); err != nil {
			ctx.Error = err
			return
		}
	}
}

// compile index page
func (b *Builder) compileIndex(ctx *Context) {
	template := "posts.html"
	if t := ctx.Theme.Template("index.html"); t != nil {
		template = "index.html"
	}

	dstFile := path.Join(ctx.DstDir, "index.html")

	ctx.Navs.Hover("home") // set hover status
	defer ctx.Navs.Reset() // remember to reset
	viewData := ctx.ViewData()
	viewData["Posts"] = ctx.indexPosts
	viewData["Pager"] = ctx.indexPager
	viewData["PostType"] = "index"

	if err := b.compileTemplate(ctx, template, viewData, dstFile); err != nil {
		ctx.Error = err
		return
	}
}

// compile template by data and write to dest file.
func (b *Builder) compileTemplate(ctx *Context, file string, viewData map[string]interface{}, destFile string) error {
	if com.IsFile(destFile) {
		ctx.Diff.Add(destFile, DiffUpdate, time.Now())
	} else {
		ctx.Diff.Add(destFile, DiffAdd, time.Now())
	}

	os.MkdirAll(path.Dir(destFile), os.ModePerm)
	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := ctx.Theme.Execute(f, file, viewData); err != nil {
		return err
	}
	return nil
}
