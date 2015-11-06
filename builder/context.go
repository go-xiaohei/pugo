package builder

import "pugo/model"

// build context, maintain parse data, posts, pages or others
type context struct {
	DstDir     string
	Posts      []*model.Post
	Pages      []*model.Page
	IndexPosts []*model.Post // temp posts for index page
	IndexPager *model.Pager
	Navs       model.Navs
}
