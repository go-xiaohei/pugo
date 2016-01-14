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
	Node          model.NodeGroup
	indexPosts    []*model.Post // temp posts for index page
	indexPager    *helper.Pager
	Data          map[string]*model.Data // custom data from other ini file

	Tags     map[string]*model.Tag
	tagPosts map[string][]*model.Post

	// site meta
	Navs model.Navs
	Meta *model.Meta

	// i18n
	I18n      *helper.I18n
	I18nGroup model.I18nGroup

	// Author
	Owner   *model.Author
	Authors model.AuthorMap

	// Comment
	Comment *model.Comment

	// Build Configuration
	Conf *model.Conf

	// Analytics
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
		"I18ns":     ctx.I18nGroup,
		"Analytics": ctx.Analytics,
		"Node":      model.NodeGroupPub(ctx.Node),
		"Lang":      ctx.Meta.Lang,
	}
	ctx.Navs.I18n(ctx.I18n)
	return m
}

// Duration returns duration of build process
func (ctx *Context) Duration() time.Duration {
	return time.Since(ctx.BeginTime)
}
