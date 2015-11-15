package builder

import "github.com/go-xiaohei/pugo-static/model"

// build context, maintain parse data, posts, pages or others
type Context struct {
	DstDir  string
	Version builderVersion

	Posts      []*model.Post
	Tags       map[string]*model.Tag
	PostPages  int
	Pages      []*model.Page
	IndexPosts []*model.Post // temp posts for index page
	IndexPager *model.Pager
	Navs       model.Navs
	Meta       *model.Meta
	Comment    *model.Comment
}

func (ctx *Context) ViewData() map[string]interface{} {
	m := map[string]interface{}{
		"Version": ctx.Version,
		"Nav":     ctx.Navs,
		"Meta":    ctx.Meta,
		"Title":   ctx.Meta.Title + " - " + ctx.Meta.Subtitle,
		"Desc":    ctx.Meta.Desc,
		"Comment": ctx.Comment,
	}
	return m
}
