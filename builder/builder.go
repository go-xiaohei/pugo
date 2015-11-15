package builder

import (
	"errors"
	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo-static/parser"
	"github.com/go-xiaohei/pugo-static/render"
	"gopkg.in/inconshreveable/log15.v2"
	"os"
	"strconv"
)

var (
	ErrSrcDirMissing = errors.New("builder-src-dir-missing")
	ErrTplDirMissing = errors.New("builder-tpl-dir-missing")
)

type (
	Builder struct {
		srcDir     string
		tplDir     string
		isBuilding bool

		renders *render.Renders
		report  *Report
		context *Context
		parser  parser.Parser
		tasks   []*BuildTask

		Error   error
		Version builderVersion
	}
	BuildTask struct {
		Name  string
		Fn    func(*Context, *Report)
		Print func(*Context) string
	}
	builderVersion struct {
		Num  string
		Date string
	}
)

func New(sourceDir, templateDir, currentTheme string, debug bool) *Builder {
	if !com.IsDir(sourceDir) {
		return &Builder{Error: ErrSrcDirMissing}
	}
	if !com.IsDir(templateDir) {
		return &Builder{Error: ErrTplDirMissing}
	}
	builder := &Builder{
		srcDir:  sourceDir,
		tplDir:  templateDir,
		parser:  parser.NewCommonParser(),
		Version: builderVersion{},
	}
	r, err := render.NewRenders(templateDir, currentTheme, debug)
	if err != nil {
		return &Builder{Error: err}
	}
	builder.renders = r
	builder.tasks = []*BuildTask{
		&BuildTask{"Meta", builder.meta, nil},
		&BuildTask{"Navs", builder.nav, func(ctx *Context) string {
			return strconv.Itoa(len(ctx.Navs))
		}},
		&BuildTask{"Comment", builder.comment, func(ctx *Context) string {
			return ctx.Comment.String()
		}},
		&BuildTask{"Posts", builder.posts, func(ctx *Context) string {
			return strconv.Itoa(len(ctx.Posts))
		}},
		&BuildTask{"Tags", builder.tags, func(ctx *Context) string {
			return strconv.Itoa(len(ctx.Tags))
		}},
		&BuildTask{"Pages", builder.pages, func(ctx *Context) string {
			return strconv.Itoa(len(ctx.Pages))
		}},
		&BuildTask{"Index", builder.index, nil},
		&BuildTask{"Feed", builder.feed, nil},
		&BuildTask{"Errors", builder.errors, nil},
		&BuildTask{"Assets", builder.assets, nil},
	}
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
		DstDir:  dest,
		Version: b.Version,
	}
	b.isBuilding = true
	for _, task := range b.tasks {
		task.Fn(ctx, r)
		if r.Error != nil {
			log15.Error("Build."+task.Name, "error", r.Error.Error())

			b.isBuilding = false
			b.report = r
			b.context = ctx
			return
		}
		if task.Print != nil {
			log15.Debug("Build."+task.Name+"."+task.Print(ctx), "duration", r.Duration())
		} else {
			log15.Debug("Build."+task.Name, "duration", r.Duration())
		}
	}
	log15.Info("Build.Finish", "duration", r.Duration(), "error", r.Error)
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
