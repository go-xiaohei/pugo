package builder

import (
	"time"

	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/go-xiaohei/pugo/app/model"
	"github.com/go-xiaohei/pugo/app/render"
)

// Context maintains parse data, posts, pages or others
type Context struct {
	Theme     *render.Theme
	DstDir    string // read output destination
	Version   string
	BeginTime time.Time
	Diff      *Diff
	Error     error

	Posts         []*model.Post
	PostPageCount int
	Pages         []*model.Page
	PageNodes     model.PageNodeGroup
	indexPosts    []*model.Post // temp posts for index page
	indexPager    *helper.Pager

	Tags     map[string]*model.Tag
	tagPosts map[string][]*model.Post

	Navs      model.Navs
	Meta      *model.Meta
	I18n      *helper.I18n // use i18n tool
	Owner     *model.Author
	Authors   model.AuthorMap
	Comment   *model.Comment
	Conf      *model.Conf
	Analytics *model.Analytics

	staticPath string
	mediaPath  string
}

// NewContext creates new build context with destination directory and version string
func NewContext(dest, ver string) *Context {
	ctx := &Context{
		DstDir:    dest,
		Version:   ver,
		BeginTime: time.Now(),
		Diff:      newDiff(),
	}
	return ctx
}

// ViewData returns global view data for template compilation
func (ctx *Context) ViewData() map[string]interface{} {
	m := map[string]interface{}{
		"Version":   ctx.Version,
		"Nav":       ctx.Navs,
		"Meta":      ctx.Meta,
		"Title":     ctx.Meta.Title + " - " + ctx.Meta.Subtitle,
		"Desc":      ctx.Meta.Desc,
		"Comment":   ctx.Comment,
		"Root":      ctx.Meta.Base,
		"Owner":     ctx.Owner,
		"I18n":      ctx.I18n,
		"Analytics": ctx.Analytics,
		"Node":      ctx.PageNodes,
	}
	return m
}

// Duration returns duration of build process
func (ctx *Context) Duration() time.Duration {
	return time.Since(ctx.BeginTime)
}
