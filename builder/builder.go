package builder

import (
	"errors"
	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo-static/parser"
	"github.com/go-xiaohei/pugo-static/render"
	"gopkg.in/inconshreveable/log15.v2"
	"os"
)

var (
	ErrSrcDirMissing = errors.New("builder-src-dir-missing")
	ErrTplDirMissing = errors.New("builder-tpl-dir-missing")
)

type Builder struct {
	srcDir     string
	tplDir     string
	isBuilding bool

	renders *render.Renders
	report  *Report
	context *Context
	parser  parser.Parser

	Error error
}

func New(sourceDir, templateDir, currentTheme string, debug bool) *Builder {
	if !com.IsDir(sourceDir) {
		return &Builder{Error: ErrSrcDirMissing}
	}
	if !com.IsDir(templateDir) {
		return &Builder{Error: ErrTplDirMissing}
	}
	builder := &Builder{
		srcDir: sourceDir,
		tplDir: templateDir,
		parser: parser.NewCommonParser(),
	}
	r, err := render.NewRenders(templateDir, currentTheme, debug)
	if err != nil {
		return &Builder{Error: err}
	}
	builder.renders = r
	return builder
}

func (b *Builder) Renders() *render.Renders {
	return b.renders
}

func (b *Builder) Build(dest string) {
	// if on build, do not again
	if b.isBuilding {
		return
	}
	log15.Debug("Build.Start")
	r := newReport(dest)
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		r.Error = err
		b.report = r
		return
	}
	ctx := &Context{
		DstDir: dest,
	}
	b.isBuilding = true
	b.meta(ctx, r)
	log15.Debug("Build.Meta", "duration", r.Duration())
	b.nav(ctx, r)
	log15.Debug("Build.Navs", "navs", len(ctx.Navs), "duration", r.Duration())
	b.comment(ctx, r)
	log15.Debug("Build.Comment", "comment", ctx.Comment.String(), "duration", r.Duration())
	b.posts(ctx, r)
	log15.Debug("Build.Posts", "posts", len(ctx.Posts), "duration", r.Duration())
	b.tags(ctx, r)
	log15.Debug("Build.Tags", "tags", len(ctx.Tags), "duration", r.Duration())
	b.pages(ctx, r)
	log15.Debug("Build.Pages", "pages", len(ctx.Pages), "duration", r.Duration())
	b.index(ctx, r)
	log15.Debug("Build.Index", "duration", r.Duration())
	b.feed(ctx, r)
	log15.Debug("Build.Feed", "duration", r.Duration())
	b.errors(ctx, r)
	log15.Debug("Build.Errors", "duration", r.Duration())
	if r.Error != nil {
		log15.Error("Build.Error", "error", r.Error.Error())
	} else {
		log15.Info("Build.Finish", "duration", r.Duration(), "error", r.Error)
	}
	b.isBuilding = false
	b.report = r
	b.context = ctx
}

func (b *Builder) IsBuilding() bool {
	return b.isBuilding
}

func (b *Builder) Report() *Report {
	return b.report
}

func (b *Builder) Context() *Context {
	return b.context
}
