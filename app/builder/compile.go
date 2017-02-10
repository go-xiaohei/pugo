package builder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/go-xiaohei/pugo/app/model"
	"github.com/gorilla/feeds"
	"gopkg.in/inconshreveable/log15.v2"
)

// Compile compile source to static files
func Compile(ctx *Context) {
	if ctx.Source == nil || ctx.Theme == nil {
		ctx.Err = fmt.Errorf("need sources data and theme to compile")
		return
	}

	// init worker in this progress
	w := helper.NewWorker(0)

	var reqs []helper.WorkerFunc
	reqs = append(reqs, compilePosts(ctx)...)
	reqs = append(reqs, compileIndexPage(ctx))
	reqs = append(reqs, compilePagePosts(ctx)...)
	reqs = append(reqs, compileTagPosts(ctx)...)
	reqs = append(reqs, compilePages(ctx)...)
	reqs = append(reqs, compileArchive(ctx))

	for _, fn := range reqs {
		w.AddFunc(fn)
	}
	w.RunOnce()
	for _, err := range w.Errors() {
		log15.Error("Build|%s", err.Error())
	}

	if ctx.Err = compileRSS(ctx); ctx.Err != nil {
		log15.Info("Compile|Done")
		return
	}
	if ctx.Err = compileSitemap(ctx); ctx.Err != nil {
		log15.Info("Compile|Done")
		return
	}
	log15.Info("Compile|Done")
}

func compilePosts(ctx *Context) []helper.WorkerFunc {
	posts := ctx.Source.Posts
	if len(posts) == 0 {
		log15.Warn("NoPosts")
		return nil
	}
	var fns []helper.WorkerFunc
	for _, post := range posts {
		p2 := post
		fn := func() error {
			viewData := ctx.View()
			viewData["Title"] = p2.Title + " - " + ctx.Source.Meta.Title
			viewData["Desc"] = p2.Desc
			viewData["Post"] = p2
			viewData["PermaKey"] = p2.Slug
			viewData["PostType"] = model.TreePost
			viewData["Hover"] = model.TreePost
			viewData["URL"] = p2.URL()
			err := compile(ctx, "post.html", viewData, p2.DestURL())
			if err != nil {
				err = fmt.Errorf("%s|%s", p2.SourceURL(), err.Error())
			}
			return err
		}
		fns = append(fns, fn)
	}
	return fns
}

func compilePagePosts(ctx *Context) []helper.WorkerFunc {
	var fns []helper.WorkerFunc
	lists := ctx.Source.PagePosts
	for page := range lists {
		pp := lists[page]
		pageKey := fmt.Sprintf("post-page-%d", pp.Pager.Current)
		fn := func() error {
			viewData := ctx.View()
			viewData["Title"] = fmt.Sprintf("Page %d - %s", pp.Pager.Current, ctx.Source.Meta.Title)
			viewData["Posts"] = pp.Posts
			viewData["Pager"] = pp.Pager
			viewData["PostType"] = model.TreePostList
			viewData["PermaKey"] = pageKey
			viewData["Hover"] = model.TreePostList
			viewData["URL"] = pp.URL
			err := compile(ctx, "posts.html", viewData, pp.DestURL())
			if err != nil {
				err = fmt.Errorf("%s|%s", pageKey, err.Error())
			}
			return err
		}
		fns = append(fns, fn)
	}
	return fns
}

func compileIndexPage(ctx *Context) helper.WorkerFunc {
	pp := ctx.Source.IndexPosts
	fn := func() error {
		template := "index.html"
		if ctx.Theme.Template(template) == nil {
			template = "posts.html"
		}
		viewData := ctx.View()
		viewData["Posts"] = pp.Posts
		viewData["Pager"] = pp.Pager
		viewData["PostType"] = model.TreeIndex
		viewData["PermaKey"] = model.TreeIndex
		viewData["Hover"] = model.TreeIndex
		viewData["URL"] = path.Join(ctx.Source.Meta.Path, "index.html")
		err := compile(ctx, template, viewData, pp.DestURL())
		if err != nil {
			err = fmt.Errorf("index.html|%s", err.Error())
		}
		return err
	}
	return fn
}

func compileTagPosts(ctx *Context) []helper.WorkerFunc {
	var fns []helper.WorkerFunc
	lists := ctx.Source.TagPosts
	for t := range lists {
		tp := lists[t]
		fn := func() error {
			pageURL := path.Join(ctx.Source.Meta.Path, tp.Tag.URL)
			pageKey := fmt.Sprintf("post-tag-%s", t)
			viewData := ctx.View()
			viewData["Title"] = fmt.Sprintf("%s - %s", t, ctx.Source.Meta.Title)
			viewData["Posts"] = tp.Posts
			viewData["Tag"] = tp.Tag
			viewData["PostType"] = model.TreePostTag
			viewData["PermaKey"] = pageKey
			viewData["Hover"] = model.TreePostTag
			viewData["URL"] = pageURL
			err := compile(ctx, "posts.html", viewData, tp.DestURL())
			if err != nil {
				err = fmt.Errorf("%s|%s", pageKey, err.Error())
			}
			return err
		}
		fns = append(fns, fn)
	}
	return fns
}

func compilePages(ctx *Context) []helper.WorkerFunc {
	pages := ctx.Source.Pages
	if len(pages) == 0 {
		log15.Warn("MoPages")
		return nil
	}
	var fns []helper.WorkerFunc
	for _, page := range pages {
		p := page
		fn := func() error {
			if p.Node {
				return nil
			}
			viewData := ctx.View()
			viewData["Title"] = p.Title + " - " + ctx.Source.Meta.Title
			viewData["Desc"] = p.Desc
			viewData["Page"] = p
			viewData["PermaKey"] = p.Slug
			viewData["PostType"] = model.TreePage
			viewData["Hover"] = p.NavHover
			viewData["URL"] = p.URL()
			if p.Lang != "" {
				viewData["Lang"] = p.Lang
				if i18n, ok := ctx.Source.I18n[p.Lang]; ok {
					viewData["I18n"] = i18n
				}
			}
			tpl := "page.html"
			if p.Template != "" {
				tpl = p.Template
			}
			err := compile(ctx, tpl, viewData, p.DestURL())
			if err != nil {
				err = fmt.Errorf("%s|%s", p.SourceURL(), err.Error())
			}
			return err
		}
		fns = append(fns, fn)
	}
	return fns
}

func compileArchive(ctx *Context) helper.WorkerFunc {
	archive := ctx.Source.Archive
	return func() error {
		viewData := ctx.View()
		viewData["Title"] = fmt.Sprintf("Archive - %s", ctx.Source.Meta.Title)
		viewData["Archives"] = archive.Data
		viewData["PostType"] = model.TreeArchive
		viewData["PermaKey"] = "archive"
		viewData["Hover"] = "archive"
		viewData["URL"] = path.Join(ctx.Source.Meta.Path, "archive")
		err := compile(ctx, "archive.html", viewData, archive.DestURL())
		if err != nil {
			err = fmt.Errorf("archive.html|%s", err.Error())
		}
		return err
	}
}

func compile(ctx *Context, file string, viewData map[string]interface{}, destFile string) error {
	os.MkdirAll(filepath.Dir(destFile), os.ModePerm)
	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := ctx.Theme.Execute(f, file, viewData); err != nil {
		return err
	}
	ctx.Sync.SetSynced(destFile)
	log15.Debug("Build|%s", filepath.ToSlash(destFile))
	atomic.AddInt32(&ctx.counter, 1)
	return nil
}

func compileRSS(ctx *Context) error {
	// todo : should compile RSS if no posts ?
	toDir := ctx.DstDir()
	now := time.Now()
	feed := &feeds.Feed{
		Title:       ctx.Source.Meta.Title,
		Link:        &feeds.Link{Href: ctx.Source.Meta.Root},
		Description: ctx.Source.Meta.Desc,
		Created:     now,
	}
	if ctx.Source.Owner != nil {
		feed.Author = &feeds.Author{
			Name:  ctx.Source.Owner.Nick,
			Email: ctx.Source.Owner.Email,
		}
	}
	var item *feeds.Item
	for _, p := range ctx.Source.Posts {
		item = &feeds.Item{
			Title:       p.Title,
			Link:        &feeds.Link{Href: ctx.Source.Meta.DomainURL(p.URL())},
			Description: string(p.Content()),
			Created:     p.Created(),
			Updated:     p.Updated(),
		}
		if p.Author != nil {
			item.Author = &feeds.Author{
				Name:  p.Author.Nick,
				Email: p.Author.Email,
			}
		}
		feed.Items = append(feed.Items, item)
	}

	dstFile := path.Join(toDir, ctx.Source.Meta.Path, "feed.xml")
	os.MkdirAll(filepath.Dir(dstFile), os.ModePerm)
	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = feed.WriteRss(f); err != nil {
		return err
	}
	ctx.Sync.SetSynced(dstFile)
	log15.Debug("Build|%s", dstFile)
	atomic.AddInt32(&ctx.counter, 1)
	return nil
}

func compileSitemap(ctx *Context) error {
	toDir := ctx.DstDir()
	now := time.Now()
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	buf.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	buf.WriteString("<url>")
	fmt.Fprintf(&buf, "<loc>%s</loc>", ctx.Source.Meta.Root)
	fmt.Fprintf(&buf, "<lastmod>%s</lastmod>", now.Format(time.RFC3339))
	buf.WriteString("<changefreq>daily</changefreq>")
	buf.WriteString("<priority>1.0</priority>")
	buf.WriteString("</url>")

	for _, p := range ctx.Source.Pages {
		buf.WriteString("<url>")
		fmt.Fprintf(&buf, "<loc>%s</loc>", ctx.Source.Meta.DomainURL(p.URL()))
		fmt.Fprintf(&buf, "<lastmod>%s</lastmod>", p.Created().Format(time.RFC3339))
		buf.WriteString("<changefreq>weekly</changefreq>")
		buf.WriteString("<priority>0.5</priority>")
		buf.WriteString("</url>")
	}

	for _, p := range ctx.Source.Posts {
		buf.WriteString("<url>")
		fmt.Fprintf(&buf, "<loc>%s</loc>", ctx.Source.Meta.DomainURL(p.URL()))
		fmt.Fprintf(&buf, "<lastmod>%s</lastmod>", p.Created().Format(time.RFC3339))
		buf.WriteString("<changefreq>daily</changefreq>")
		buf.WriteString("<priority>0.6</priority>")
		buf.WriteString("</url>")
	}
	buf.WriteString("<url>")
	fmt.Fprintf(&buf, "<loc>%s</loc>", ctx.Source.Meta.DomainURL("archive.html"))
	fmt.Fprintf(&buf, "<lastmod>%s</lastmod>", now.Format(time.RFC3339))
	buf.WriteString("<changefreq>daily</changefreq>")
	buf.WriteString("<priority>0.6</priority>")
	buf.WriteString("</url>")

	for i := 1; i <= ctx.Source.PostPage; i++ {
		buf.WriteString("<url>")
		fmt.Fprintf(&buf, "<loc>%s</loc>", ctx.Source.Meta.DomainURL(fmt.Sprintf("post/%d.html", i)))
		fmt.Fprintf(&buf, "<lastmod>%s</lastmod>", now.Format(time.RFC3339))
		buf.WriteString("<changefreq>daily</changefreq>")
		buf.WriteString("<priority>0.6</priority>")
		buf.WriteString("</url>")
	}

	for _, t := range ctx.Source.Tags {
		buf.WriteString("<url>")
		fmt.Fprintf(&buf, "<loc>%s</loc>", ctx.Source.Meta.DomainURL(t.URL))
		fmt.Fprintf(&buf, "<lastmod>%s</lastmod>", now.Format(time.RFC3339))
		buf.WriteString("<changefreq>weekly</changefreq>")
		buf.WriteString("<priority>0.5</priority>")
		buf.WriteString("</url>")
	}

	buf.WriteString("</urlset>")
	dstFile := path.Join(toDir, ctx.Source.Meta.Path, "sitemap.xml")
	os.MkdirAll(path.Dir(dstFile), os.ModePerm)
	if err := ioutil.WriteFile(dstFile, buf.Bytes(), os.ModePerm); err != nil {
		return err
	}
	ctx.Sync.SetSynced(dstFile)
	log15.Debug("Build|%s", dstFile)
	atomic.AddInt32(&ctx.counter, 1)
	return nil
}
