package builder

import (
	"bytes"
	"context"
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
	worker := helper.NewGoWorker()
	worker.Start()
	// add worker result handler
	worker.Receive(func(rs *helper.GoWorkerResult) {
		if rs.Error != nil {
			if p, ok := rs.Ctx.Value("post").(*model.Post); ok {
				log15.Error("Build|%s|%s", p.SourceURL(), rs.Error.Error())
			}
			if p, ok := rs.Ctx.Value("page").(*model.Page); ok {
				log15.Error("Build|%s|%s", p.SourceURL(), rs.Error.Error())
			}
		}
	})
	var reqs []*helper.GoWorkerRequest
	reqs = append(reqs, compilePosts2(ctx)...)
	reqs = append(reqs, compileIndexPage(ctx))
	reqs = append(reqs, compilePagePosts(ctx)...)
	reqs = append(reqs, compileTagPosts(ctx)...)
	reqs = append(reqs, compilePages2(ctx)...)
	reqs = append(reqs, compileArchive(ctx))
	for _, r := range reqs {
		worker.Send(r)
	}
	worker.WaitStop()
	ctx.Err = compileXML(ctx)
}

func compilePosts2(ctx *Context) []*helper.GoWorkerRequest {
	posts := ctx.Source.Posts
	if len(posts) == 0 {
		log15.Warn("NoPosts")
		return nil
	}
	var reqs []*helper.GoWorkerRequest

	for _, post := range posts {
		p2 := post
		c := context.WithValue(context.Background(), "post", p2)
		req := &helper.GoWorkerRequest{
			Ctx: c,
			Action: func(c context.Context) (context.Context, error) {
				viewData := ctx.View()
				viewData["Title"] = p2.Title + " - " + ctx.Source.Meta.Title
				viewData["Desc"] = p2.Desc
				viewData["Post"] = p2
				viewData["PermaKey"] = p2.Slug
				viewData["PostType"] = model.TreePost
				viewData["Hover"] = model.TreePost
				viewData["URL"] = p2.URL()
				return c, compile(ctx, "post.html", viewData, p2.DestURL())
			},
		}
		reqs = append(reqs, req)
	}
	return reqs
}

func compilePagePosts(ctx *Context) []*helper.GoWorkerRequest {
	var reqs []*helper.GoWorkerRequest
	lists := ctx.Source.PagePosts
	for page := range lists {
		pp := lists[page]
		pageKey := fmt.Sprintf("post-page-%d", pp.Pager.Current)
		c := context.WithValue(context.Background(), "post-page", pp.Posts)
		c = context.WithValue(c, "page", pp.Pager)
		req := &helper.GoWorkerRequest{
			Ctx: c,
			Action: func(c context.Context) (context.Context, error) {
				viewData := ctx.View()
				viewData["Title"] = fmt.Sprintf("Page %d - %s", pp.Pager.Current, ctx.Source.Meta.Title)
				viewData["Posts"] = pp.Posts
				viewData["Pager"] = pp.Pager
				viewData["PostType"] = model.TreePostList
				viewData["PermaKey"] = pageKey
				viewData["Hover"] = model.TreePostList
				viewData["URL"] = pp.URL
				return c, compile(ctx, "posts.html", viewData, pp.DestURL())
			},
		}
		reqs = append(reqs, req)
	}
	return reqs
}

func compileIndexPage(ctx *Context) *helper.GoWorkerRequest {
	pp := ctx.Source.IndexPosts
	c := context.WithValue(context.Background(), "post-page", pp.Posts)
	c = context.WithValue(c, "page", pp.Pager)
	return &helper.GoWorkerRequest{
		Ctx: c,
		Action: func(c context.Context) (context.Context, error) {
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
			return c, compile(ctx, template, viewData, pp.DestURL())
		},
	}
}

func compileTagPosts(ctx *Context) []*helper.GoWorkerRequest {
	var reqs []*helper.GoWorkerRequest
	lists := ctx.Source.TagPosts
	for t := range lists {
		tp := lists[t]
		c := context.WithValue(context.Background(), "post-tag", tp.Posts)
		c = context.WithValue(c, "tag", tp.Tag)
		req := &helper.GoWorkerRequest{
			Ctx: c,
			Action: func(c context.Context) (context.Context, error) {
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
				return c, compile(ctx, "posts.html", viewData, tp.DestURL())
			},
		}
		reqs = append(reqs, req)
	}
	return reqs
}

func compilePages2(ctx *Context) []*helper.GoWorkerRequest {
	pages := ctx.Source.Pages
	if len(pages) == 0 {
		log15.Warn("MoPages")
		return nil
	}
	var reqs []*helper.GoWorkerRequest
	for _, page := range pages {
		p := page
		c := context.WithValue(context.Background(), "page", p)
		req := &helper.GoWorkerRequest{
			Ctx: c,
			Action: func(c context.Context) (context.Context, error) {
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
				return c, compile(ctx, tpl, viewData, p.DestURL())
			},
		}
		reqs = append(reqs, req)
	}
	return reqs
}

func compileArchive(ctx *Context) *helper.GoWorkerRequest {
	archive := ctx.Source.Archive
	c := context.WithValue(context.Background(), "Archives", archive.Data)
	return &helper.GoWorkerRequest{
		Ctx: c,
		Action: func(c context.Context) (context.Context, error) {
			viewData := ctx.View()
			viewData["Title"] = fmt.Sprintf("Archive - %s", ctx.Source.Meta.Title)
			viewData["Archives"] = archive.Data
			viewData["PostType"] = model.TreeArchive
			viewData["PermaKey"] = "archive"
			viewData["Hover"] = "archive"
			viewData["URL"] = path.Join(ctx.Source.Meta.Path, "archive")
			return c, compile(ctx, "archive.html", viewData, archive.DestURL())
		},
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
	ctx.Files.Add(destFile, 0, ctx.time, model.FileCompiled, model.OpCompiled)
	log15.Debug("Build|%s", filepath.ToSlash(destFile))
	atomic.AddInt64(&ctx.counter, 1)
	return nil
}

func compileXML(ctx *Context) error {
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
	ctx.Files.Add(dstFile, 0, ctx.time, model.FileCompiled, model.OpCompiled)
	log15.Debug("Build|%s", dstFile)
	atomic.AddInt64(&ctx.counter, 1)

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
	dstFile = path.Join(toDir, ctx.Source.Meta.Path, "sitemap.xml")
	os.MkdirAll(path.Dir(dstFile), os.ModePerm)
	if err = ioutil.WriteFile(dstFile, buf.Bytes(), os.ModePerm); err != nil {
		return err
	}
	ctx.Files.Add(dstFile, 0, ctx.time, model.FileCompiled, model.OpCompiled)
	log15.Debug("Build|%s", dstFile)
	atomic.AddInt64(&ctx.counter, 1)
	return nil
}
