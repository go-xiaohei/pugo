package builder

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/parser"
	"github.com/go-xiaohei/pugo/app/render"
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
		isWatching bool
		opt        *BuildOption

		render  *render.Render
		context *Context
		parsers []parser.Parser
		tasks   []*BuildTask

		Error   error
		Version builderVersion
		Count   uint32 // build count
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

		beforeHook []BuildHook
		afterHook  []BuildHook
	}
	// hook func to build process
	BuildHook func(b *Builder, ctx *Context) error
)

func (opt *BuildOption) Before(fn BuildHook) {
	opt.beforeHook = append(opt.beforeHook, fn)
}

func (opt *BuildOption) After(fn BuildHook) {
	opt.afterHook = append(opt.afterHook, fn)
}

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
	b.isBuilding = true

	ctx := &Context{
		DstDir:    dest,
		Version:   b.Version,
		BeginTime: time.Now(),
		Diff:      newDiff(),
	}
	// before hook
	if len(b.opt.beforeHook) > 0 {
		for _, fn := range b.opt.beforeHook {
			if err := fn(b, ctx); err != nil {
				log15.Error("Build.Before", "error", err.Error())
				ctx.Error = err
				b.isBuilding = false
				b.context = ctx
				return
			}
		}
	}

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

	b.isBuilding = false
	atomic.AddUint32(&b.Count, 1)
	log15.Info("Build.Finish", "duration", ctx.Duration(), "count", b.Count)

	// after hook
	if len(b.opt.afterHook) > 0 {
		for _, fn := range b.opt.afterHook {
			if err := fn(b, ctx); err != nil {
				log15.Error("Build.After", "error", err.Error())
				if ctx.Error == nil {
					ctx.Error = err
				}
				b.context = ctx
				return
			}
		}
	}
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

// is builder watching changes
func (b *Builder) IsWatching() bool {
	return b.isWatching
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
