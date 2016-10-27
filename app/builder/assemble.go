package builder

import (
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/inconshreveable/log15.v2"

	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/go-xiaohei/pugo/app/model"
)

// AssembleSource assemble some extra data in Source,
// such as page nodes, i18n status.
// it need be used after posts and pages are loaded
func AssembleSource(ctx *Context) {
	if ctx.Source == nil || ctx.Theme == nil {
		ctx.Err = fmt.Errorf("need sources data and theme to assemble")
		return
	}

	ctx.Source.Nav.SetPrefix(ctx.Source.Meta.Path)
	ctx.Source.Tags = make(map[string]*model.Tag)
	ctx.Source.TagPosts = make(map[string]*model.TagPosts)
	ctx.Source.PagePosts = make(map[int]*model.PagerPosts)

	r, hr := newReplacer(ctx.Source.Meta.Path), newReplacerInHTML(ctx.Source.Meta.Path)
	ctx.Source.Meta.Cover = r.Replace(ctx.Source.Meta.Cover)
	for _, a := range ctx.Source.Authors {
		a.Avatar = r.Replace(a.Avatar)
	}

	// fill post data
	for _, p := range ctx.Source.Posts {
		if ctx.Source.Meta.Path != "" && ctx.Source.Meta.Path != "/" {
			p.SetURL(path.Join(ctx.Source.Meta.Path, p.URL()))
		}
		p.SetDestURL(filepath.Join(ctx.DstDir(), p.URL()))
		p.SetPlaceholder(r, hr)
		ctx.Tree.Add(p.DestURL(), p.Title, model.TreePost, 0)
		if p.Author == nil {
			p.Author = ctx.Source.Authors[p.AuthorName]
		}
		for _, t := range p.Tags {
			ctx.Source.Tags[t.Name] = t
			if ctx.Source.TagPosts[t.Name] == nil {
				ctx.Source.TagPosts[t.Name] = &model.TagPosts{
					Posts: []*model.Post{p},
					Tag:   t,
				}
			} else {
				ctx.Source.TagPosts[t.Name].Posts = append(ctx.Source.TagPosts[t.Name].Posts, p)
			}
		}
	}

	// fill page data
	for _, p := range ctx.Source.Pages {
		if ctx.Source.Meta.Path != "" && ctx.Source.Meta.Path != "/" {
			p.SetURL(path.Join(ctx.Source.Meta.Path, p.URL()))
		}
		p.SetDestURL(filepath.Join(ctx.DstDir(), p.URL()))
		p.SetPlaceholder(hr)
		treeType := model.TreePage
		if p.Node {
			treeType = model.TreePageNode
		}
		ctx.Tree.Add(p.DestURL(), p.Title, treeType, p.Sort)
		if p.Author == nil {
			p.Author = ctx.Source.Authors[p.AuthorName]
		}
	}

	// prepare tag posts
	for _, tp := range ctx.Source.TagPosts {
		sort.Sort(model.Posts(tp.Posts))
		tp.SetDestURL(path.Join(ctx.DstDir(),
			ctx.Source.Meta.Path, tp.Tag.URL))
		ctx.Tree.Add(tp.DestURL(), "", model.TreePostTag, 0)
	}

	// prepare archives
	archives := model.NewArchive(ctx.Source.Posts)
	archives.SetDestURL(filepath.Join(ctx.DstDir(), archives.DestURL()))
	ctx.Source.Archive = archives
	ctx.Tree.Add(archives.DestURL(), "Archive", model.TreeArchive, 0)

	// prepare paged posts
	var (
		cursor = helper.NewPagerCursor(4, len(ctx.Source.Posts))
		page   = 1
		layout = "posts/%d"
	)
	for {
		pager := cursor.Page(page)
		if pager == nil {
			ctx.Source.PostPage = page - 1
			break
		}
		currentPosts := ctx.Source.Posts[pager.Begin:pager.End]
		pager.SetLayout(path.Join(ctx.Source.Meta.Path, "/"+layout+".html"))
		pageURL := path.Join(ctx.Source.Meta.Path, fmt.Sprintf(layout+".html", pager.Current))
		pp := &model.PagerPosts{
			Posts: currentPosts,
			Pager: pager,
			URL:   pageURL,
		}
		pp.SetDestURL(path.Join(ctx.DstDir(), pageURL))
		ctx.Source.PagePosts[pager.Current] = pp
		ctx.Tree.Add(pp.DestURL(), "", model.TreePostList, 0)
		if pager.Current == 1 {
			// use new object, not pp
			pp2 := model.PagerPosts{
				Posts: currentPosts,
				Pager: pager,
			}
			pp2.SetDestURL(path.Join(ctx.DstDir(), "index.html"))
			ctx.Source.IndexPosts = pp2
			ctx.Tree.Add(path.Join(ctx.DstDir(), "index.html"), "Home", model.TreeIndex, 0)
		}
		page++
	}

	ctx.Tree.Add(path.Join(ctx.DstDir(), ctx.Source.Meta.Path, "feed.xml"), "Feed", model.TreeXML, 0)
	ctx.Tree.Add(path.Join(ctx.DstDir(), ctx.Source.Meta.Path, "sitemap.xml"), "Sitemap", model.TreeXML, 0)

	if ctx.Err = ctx.Theme.Load(); ctx.Err != nil {
		return
	}

	log15.Info("Assemble|Done")
}

func newReplacer(static string) *strings.Replacer {
	p := path.Join(static, "media")
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return strings.NewReplacer(
		"@media", p,
	)
}

func newReplacerInHTML(static string) *strings.Replacer {
	p := path.Join(static, "media")
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return strings.NewReplacer(
		`src="@media`, fmt.Sprintf(`src="%s`, p),
		`href="@media`, fmt.Sprintf(`src="%s`, p),
	)
}
