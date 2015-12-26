package builder

import (
	"time"

	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/go-xiaohei/pugo/app/model"
	"github.com/go-xiaohei/pugo/app/render"
)

// build context, maintain parse data, posts, pages or others
type Context struct {
	Theme     *render.Theme
	DstDir    string // read output destination
	Version   builderVersion
	BeginTime time.Time
	Diff      *Diff
	Error     error

	Posts         []*model.Post
	PostPageCount int
	Pages         []*model.Page
	indexPosts    []*model.Post // temp posts for index page
	indexPager    *helper.Pager

	Tags     map[string]*model.Tag
	tagPosts map[string][]*model.Post

	Navs    model.Navs
	Meta    *model.Meta
	Authors model.AuthorMap
	Comment *model.Comment
	Conf    *model.Conf

	copy2StaticDir string
	copy2MediaDir  string
}

// return global view data for template compilation
func (ctx *Context) ViewData() map[string]interface{} {
	m := map[string]interface{}{
		"Version": ctx.Version,
		"Nav":     ctx.Navs,
		"Meta":    ctx.Meta,
		"Title":   ctx.Meta.Title + " - " + ctx.Meta.Subtitle,
		"Desc":    ctx.Meta.Desc,
		"Comment": ctx.Comment,
		"Root":    ctx.Meta.Base,
	}
	return m
}

// return duration
func (ctx *Context) Duration() time.Duration {
	return time.Since(ctx.BeginTime)
}
