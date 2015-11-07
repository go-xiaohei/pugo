package builder

import "pugo/model"

// build context, maintain parse data, posts, pages or others
type context struct {
	DstDir     string
	Posts      []*model.Post
	Tags       map[string]*model.Tag
	PostPages  int
	Pages      []*model.Page
	IndexPosts []*model.Post // temp posts for index page
	IndexPager *model.Pager
	Navs       model.Navs
	Meta       *model.Meta
}

func (ctx *context) viewData() map[string]interface{} {
	m := map[string]interface{}{
		"Nav":   ctx.Navs,
		"Meta":  ctx.Meta,
		"Title": ctx.Meta.Title + " - " + ctx.Meta.Subtitle,
		"Desc":  ctx.Meta.Desc,
	}
	return m
}
