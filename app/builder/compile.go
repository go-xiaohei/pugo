package builder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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
	var destDir = ""
	if destDir, ctx.Err = toDir(ctx.To); ctx.Err != nil {
		return
	}
	if ctx.Err = compilePosts(ctx, destDir); ctx.Err != nil {
		return
	}
	if ctx.Err = compilePages(ctx, destDir); ctx.Err != nil {
		return
	}
	if ctx.Err = compileXML(ctx, destDir); ctx.Err != nil {
		return
	}
	ctx.Source.Tree.Print("")
}

func compilePosts(ctx *Context, toDir string) error {
	var (
		viewData map[string]interface{}
		dstFile  string
		err      error
	)
	for _, p := range ctx.Source.Posts {
		dstFile = path.Join(toDir, p.URL())
		if path.Ext(dstFile) == "" {
			dstFile += ".html"
		}

		viewData = ctx.View()
		viewData["Title"] = p.Title + " - " + ctx.Source.Meta.Title
		viewData["Desc"] = p.Desc
		viewData["Post"] = p
		viewData["Permalink"] = p.Permalink()
		viewData["PermaKey"] = p.Slug
		viewData["PostType"] = model.TreePost
		viewData["Hover"] = model.TreePost

		if err = compile(ctx, "post.html", viewData, dstFile); err != nil {
			return err
		}

		ctx.Source.Tree.Add(p.TreeURL(), model.TreePost, 0)
	}

	// build posts
	var (
		currentPosts []*model.Post
		cursor       = helper.NewPagerCursor(4, len(ctx.Source.Posts))
		page         = 1
		layout       = "posts/%d"
		pageKey      = ""
	)
	for {
		pager := cursor.Page(page)
		if pager == nil {
			ctx.Source.PostPage = page - 1
			break
		}

		currentPosts = ctx.Source.Posts[pager.Begin:pager.End]
		pager.SetLayout(path.Join(ctx.Source.Meta.Path, "/"+layout+".html"))

		dstFile = path.Join(toDir, ctx.Source.Meta.Path, fmt.Sprintf(layout+".html", pager.Current))

		pageKey = fmt.Sprintf("post-page-%d", pager.Current)
		viewData = ctx.View()
		viewData["Title"] = fmt.Sprintf("Page %d - %s", pager.Current, ctx.Source.Meta.Title)
		viewData["Posts"] = currentPosts
		viewData["Pager"] = pager
		viewData["PostType"] = model.TreePostList
		viewData["PermaKey"] = pageKey
		viewData["Hover"] = model.TreePostList

		if err = compile(ctx, "posts.html", viewData, dstFile); err != nil {
			return err
		}

		ctx.Source.Tree.Add(fmt.Sprintf(layout, pager.Current), model.TreePostList, 0)

		if page == 1 {

			template := "index.html"
			if ctx.Theme.Template(template) == nil {
				template = "posts.html"
			}
			viewData = ctx.View()
			viewData["Posts"] = currentPosts
			viewData["Pager"] = pager
			viewData["PostType"] = model.TreeIndex
			viewData["PermaKey"] = model.TreeIndex
			viewData["Hover"] = model.TreeIndex

			dstFile = path.Join(toDir, ctx.Source.Meta.Path, "index.html")

			if err = compile(ctx, "posts.html", viewData, dstFile); err != nil {
				return err
			}

			ctx.Source.Tree.Add("index.html", model.TreeIndex, 0)
		}
		page++
	}

	// build archive
	dstFile = path.Join(toDir, ctx.Source.Meta.Path, "archive.html")
	viewData = ctx.View()
	viewData["Title"] = fmt.Sprintf("Archive - %s", ctx.Source.Meta.Title)
	viewData["Archives"] = ctx.Source.Archive
	viewData["PostType"] = model.TreeArchive
	viewData["PermaKey"] = model.TreeArchive
	viewData["Hover"] = model.TreeArchive
	if err = compile(ctx, "archive.html", viewData, dstFile); err != nil {
		return err
	}
	ctx.Source.Tree.Add("archive.html", model.TreeArchive, 0)

	// build tag posts
	for t, posts := range ctx.Source.tagPosts {
		dstFile = path.Join(toDir, ctx.Source.Meta.Path, ctx.Source.Tags[t].URL)
		viewData = ctx.View()
		viewData["Title"] = fmt.Sprintf("%s - %s", t, ctx.Source.Meta.Title)
		viewData["Posts"] = posts
		viewData["Tag"] = ctx.Source.Tags[t]
		viewData["PostType"] = model.TreePostTag
		viewData["PermaKey"] = pageKey
		viewData["Hover"] = model.TreePostTag
		if err = compile(ctx, "posts.html", viewData, dstFile); err != nil {
			return err
		}
		ctx.Source.Tree.Add(path.Join(ctx.Source.Meta.Path, ctx.Source.Tags[t].URL), model.TreePostTag, 0)
	}
	return nil
}

func compilePages(ctx *Context, toDir string) error {
	var (
		viewData map[string]interface{}
		dstFile  string
		err      error
	)
	for _, p := range ctx.Source.Pages {
		dstFile = path.Join(toDir, p.URL())
		if path.Ext(dstFile) == "" {
			dstFile += ".html"
		}

		viewData = ctx.View()
		viewData["Title"] = p.Title + " - " + ctx.Source.Meta.Title
		viewData["Desc"] = p.Desc
		viewData["Post"] = p
		viewData["Permalink"] = p.Permalink()
		viewData["PermaKey"] = p.Slug
		viewData["PostType"] = model.TreePage
		viewData["Hover"] = model.TreePage

		if err = compile(ctx, "page.html", viewData, dstFile); err != nil {
			return err
		}

		ctx.Source.Tree.Add(p.TreeURL(), model.TreePage, p.Sort)
	}
	return nil
}

func compile(ctx *Context, file string, viewData map[string]interface{}, destFile string) error {
	os.MkdirAll(path.Dir(destFile), os.ModePerm)
	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := ctx.Theme.Execute(f, file, viewData); err != nil {
		return err
	}
	log15.Debug("Build|To|%s", destFile)
	atomic.AddInt64(&ctx.counter, 1)
	return nil
}

func compileXML(ctx *Context, toDir string) error {
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
	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = feed.WriteRss(f); err != nil {
		return err
	}
	log15.Debug("Build|To|%s", dstFile)
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
	log15.Debug("Build|To|%s", dstFile)
	atomic.AddInt64(&ctx.counter, 1)
	return nil
}
