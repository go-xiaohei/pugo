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

	// do compile task in goroutine
	wg := helper.NewGoGroup("BuildStep")
	wg.Wrap("compilePosts", func() {
		ctx.Err = compilePosts(ctx, worker, ctx.dstDir)
	})
	wg.Wrap("compilePages", func() {
		ctx.Err = compilePages(ctx, worker, ctx.dstDir)
	})
	wg.Wrap("compileXML", func() {
		ctx.Err = compileXML(ctx, ctx.dstDir)
	})
	wg.Wait()

	// wait worker
	worker.WaitStop()

}

func compilePosts(ctx *Context, w *helper.GoWorker, toDir string) error {

	// compile each post
	compilePostFn := func(ctx *Context, i int) error {
		p := ctx.Source.Posts[i]
		dstFile := path.Join(toDir, p.URL())
		if path.Ext(dstFile) == "" {
			dstFile += ".html"
		}

		viewData := ctx.View()
		viewData["Title"] = p.Title + " - " + ctx.Source.Meta.Title
		viewData["Desc"] = p.Desc
		viewData["Post"] = p
		viewData["PermaKey"] = p.Slug
		viewData["PostType"] = model.TreePost
		viewData["Hover"] = model.TreePost
		viewData["URL"] = p.URL()

		if err := compile(ctx, "post.html", viewData, dstFile); err != nil {
			return err
		}

		ctx.Tree.Add(p.TreeURL(), model.TreePost, 0)
		return nil
	}
	for i := range ctx.Source.Posts {
		i2 := i
		c := context.WithValue(context.Background(), "post", ctx.Source.Posts[i2])
		w.Send(&helper.GoWorkerRequest{
			Ctx: c,
			Action: func(c context.Context) (context.Context, error) {
				return c, compilePostFn(ctx, i2)
			},
		})
	}

	// build posts
	var (
		cursor = helper.NewPagerCursor(4, len(ctx.Source.Posts))
		page   = 1
		layout = "posts/%d"
	)

	// compile post-page
	compilePostPageFn := func(ctx *Context, pager *helper.Pager, currentPosts []*model.Post) error {
		pager.SetLayout(path.Join(ctx.Source.Meta.Path, "/"+layout+".html"))
		pageURL := path.Join(ctx.Source.Meta.Path, fmt.Sprintf(layout+".html", pager.Current))
		dstFile := path.Join(toDir, pageURL)

		pageKey := fmt.Sprintf("post-page-%d", pager.Current)
		viewData := ctx.View()
		viewData["Title"] = fmt.Sprintf("Page %d - %s", pager.Current, ctx.Source.Meta.Title)
		viewData["Posts"] = currentPosts
		viewData["Pager"] = pager
		viewData["PostType"] = model.TreePostList
		viewData["PermaKey"] = pageKey
		viewData["Hover"] = model.TreePostList
		viewData["URL"] = pageURL

		if err := compile(ctx, "posts.html", viewData, dstFile); err != nil {
			return err
		}

		ctx.Tree.Add(fmt.Sprintf(layout, pager.Current), model.TreePostList, 0)

		if pager.Current == 1 {

			template := "index.html"
			if ctx.Theme.Template(template) == nil {
				template = "posts.html"
			}
			viewData := ctx.View()
			viewData["Posts"] = currentPosts
			viewData["Pager"] = pager
			viewData["PostType"] = model.TreeIndex
			viewData["PermaKey"] = model.TreeIndex
			viewData["Hover"] = model.TreeIndex
			viewData["URL"] = path.Join(ctx.Source.Meta.Path, "index.html")

			dstFile = path.Join(toDir, ctx.Source.Meta.Path, "index.html")

			if err := compile(ctx, template, viewData, dstFile); err != nil {
				return err
			}

			ctx.Tree.Add("index.html", model.TreeIndex, 0)
		}

		return nil
	}

	for {
		pager := cursor.Page(page)
		if pager == nil {
			ctx.Source.PostPage = page - 1
			break
		}
		currentPosts := ctx.Source.Posts[pager.Begin:pager.End]
		c := context.WithValue(context.Background(), "post-page", currentPosts)
		c = context.WithValue(c, "page", page)
		w.Send(&helper.GoWorkerRequest{
			Ctx: c,
			Action: func(c context.Context) (context.Context, error) {
				return c, compilePostPageFn(ctx, pager, currentPosts)
			},
		})
		page++
	}

	// build archive
	c := context.WithValue(context.Background(), "Archives", ctx.Source.Archive)
	w.Send(&helper.GoWorkerRequest{
		Ctx: c,
		Action: func(c context.Context) (context.Context, error) {
			return c, func(ctx *Context) error {
				dstFile := path.Join(toDir, ctx.Source.Meta.Path, "archive.html")
				viewData := ctx.View()
				viewData["Title"] = fmt.Sprintf("Archive - %s", ctx.Source.Meta.Title)
				viewData["Archives"] = ctx.Source.Archive
				viewData["PostType"] = model.TreeArchive
				viewData["PermaKey"] = "archive"
				viewData["Hover"] = "archive"
				viewData["URL"] = path.Join(ctx.Source.Meta.Path, "archive")
				if err := compile(ctx, "archive.html", viewData, dstFile); err != nil {
					return err
				}
				ctx.Tree.Add("archive.html", model.TreeArchive, 0)
				return nil
			}(ctx)
		},
	})

	// compile tag posts
	compilePostTagFn := func(ctx *Context, t string) error {
		posts := ctx.Source.tagPosts[t]
		pageURL := path.Join(ctx.Source.Meta.Path, ctx.Source.Tags[t].URL)
		dstFile := path.Join(toDir, pageURL)
		pageKey := fmt.Sprintf("post-tag-%s", t)
		viewData := ctx.View()
		viewData["Title"] = fmt.Sprintf("%s - %s", t, ctx.Source.Meta.Title)
		viewData["Posts"] = posts
		viewData["Tag"] = ctx.Source.Tags[t]
		viewData["PostType"] = model.TreePostTag
		viewData["PermaKey"] = pageKey
		viewData["Hover"] = model.TreePostTag
		viewData["URL"] = pageURL
		if err := compile(ctx, "posts.html", viewData, dstFile); err != nil {
			return err
		}
		ctx.Tree.Add(path.Join(ctx.Source.Meta.Path, ctx.Source.Tags[t].URL), model.TreePostTag, 0)
		return nil
	}

	// build tag posts
	for t := range ctx.Source.tagPosts {
		t2 := t
		c := context.WithValue(context.Background(), "post-tag", ctx.Source.tagPosts[t2])
		c = context.WithValue(c, "tag", ctx.Source.Tags[t2])
		w.Send(&helper.GoWorkerRequest{
			Ctx: c,
			Action: func(c context.Context) (context.Context, error) {
				return c, compilePostTagFn(ctx, t2)
			},
		})
	}
	return nil
}

func compilePages(ctx *Context, w *helper.GoWorker, toDir string) error {

	compileFn := func(ctx *Context, i int) error {
		p := ctx.Source.Pages[i]
		dstFile := filepath.Join(toDir, p.URL())
		if path.Ext(dstFile) == "" {
			dstFile += ".html"
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
		if err := compile(ctx, tpl, viewData, dstFile); err != nil {
			return err
		}

		ctx.Tree.Add(p.TreeURL(), model.TreePage, p.Sort)

		return nil
	}

	for i := range ctx.Source.Pages {
		i2 := i
		c := context.WithValue(context.Background(), "page", ctx.Source.Pages[i2])
		w.Send(&helper.GoWorkerRequest{
			Ctx: c,
			Action: func(c context.Context) (context.Context, error) {
				return c, compileFn(ctx, i2)
			},
		})
	}

	return nil
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
