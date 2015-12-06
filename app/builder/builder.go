package builder

import (
	"errors"
	"time"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo-static/app/parser"
	"github.com/go-xiaohei/pugo-static/app/render"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	ErrSrcDirMissing = errors.New("builder-src-dir-missing")
	ErrTplDirMissing = errors.New("builder-tpl-dir-missing")
)

type (
	// builder object, provides api to build and watch sources and templates
	Builder struct {
		isBuilding bool
		opt        *BuildOption

		render  *render.Render
		context *Context
		parsers []parser.Parser
		tasks   []*BuildTask

		Error   error
		Version builderVersion
	}
	// build task defines the build function to run in build process
	BuildTask struct {
		Name  string
		Fn    func(*Context)
		Print func(*Context) string
	}
	builderVersion struct {
		Num  string
		Date string
	}
	// build option to builder
	BuildOption struct {
		SrcDir   string
		TplDir   string
		MediaDir string
		Theme    string

		Version string
		VerDate string
	}
)

// New builder with option
func New(opt *BuildOption) *Builder {
	if !com.IsDir(opt.SrcDir) {
		return &Builder{Error: ErrSrcDirMissing}
	}
	if !com.IsDir(opt.TplDir) {
		return &Builder{Error: ErrTplDirMissing}
	}
	builder := &Builder{
		parsers: []parser.Parser{
			parser.NewCommonParser(),
			parser.NewMdParser(),
		},
		Version: builderVersion{
			Num:  opt.Version,
			Date: opt.VerDate,
		},
		opt: opt,
	}
	builder.render = render.New(builder.opt.TplDir)
	builder.tasks = []*BuildTask{
		{"Data", builder.ReadData, nil},
		{"Compile", builder.Compile, nil},
		{"Feed", builder.WriteFeed, nil},
		{"Copy", builder.CopyAssets, nil},
	}
	log15.Debug("Build.Source." + opt.SrcDir)
	log15.Debug("Build.Template." + opt.TplDir)
	log15.Debug("Build.Theme." + opt.Theme)
	return builder
}

// get render in builder
func (b *Builder) Render() *render.Render {
	return b.render
}

// build to dest directory
func (b *Builder) Build(dest string) {
	// if on build, do not again
	if b.isBuilding {
		return
	}
	log15.Debug("Build.Start")
	ctx := &Context{
		DstDir:       dest,
		DstOriginDir: dest,
		Version:      b.Version,
		BeginTime:    time.Now(),
	}
	b.isBuilding = true

	// run tasks
	for _, task := range b.tasks {
		task.Fn(ctx)
		if ctx.Error != nil {
			log15.Error("Build."+task.Name, "error", ctx.Error.Error())

			b.isBuilding = false
			b.context = ctx
			return
		}
		if task.Print != nil {
			log15.Debug("Build."+task.Name+"."+task.Print(ctx), "duration", ctx.Duration())
		} else {
			log15.Debug("Build."+task.Name, "duration", ctx.Duration())
		}
		b.context = ctx
	}

	log15.Info("Build.Finish", "duration", ctx.Duration())
	b.isBuilding = false
}

// get parser with mark bytes
func (b *Builder) getParser(data []byte) parser.Parser {
	for _, p := range b.parsers {
		if p.Is(data) {
			return p
		}
	}
	return nil
}

// is builder run building
func (b *Builder) IsBuilding() bool {
	return b.isBuilding
}

// get last context in builder
func (b *Builder) Context() *Context {
	return b.context
}

// get option if nil, or set option with non-nil opt.
func (b *Builder) Option(opt *BuildOption) *BuildOption {
	if opt == nil {
		return b.opt
	}
	b.opt = opt
	return nil
}
