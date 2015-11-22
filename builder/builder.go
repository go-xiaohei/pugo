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

type (
	// builder object, provides api to build and watch sources and templates
	Builder struct {
		isBuilding bool
		opt        *BuildOption

		render  *render.Render
		report  *Report
		context *Context
		parsers []parser.Parser
		tasks   []*BuildTask

		Error   error
		Version builderVersion
	}
	// build task defines the build function to run in build process
	BuildTask struct {
		Name  string
		Fn    func(*Context, *Report)
		Print func(*Context) string
	}
	builderVersion struct {
		Num  string
		Date string
	}
	// build option to builder
	BuildOption struct {
		SrcDir    string
		TplDir    string
		Theme     string
		UploadDir string

		Version string
		VerDate string

		IsDebug         bool
		IsCopyAssets    bool
		IsWatchTemplate bool
		IsSuffixed      bool
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
		&BuildTask{"Data", builder.ReadData, nil},
		&BuildTask{"Compile", builder.Compile, nil},
		&BuildTask{"Copy", builder.CopyAssets, nil},
		&BuildTask{"Feed", builder.WriteFeed, nil},
	}
	return builder
}

// get render in builder
func (b *Builder) Render() *render.Render {
	return b.render
}

// get current theme in render
func (b *Builder) theme() (*render.Theme, error) {
	return b.render.Load(b.opt.Theme)
}

// build to dest directory
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
	theme, err := b.theme()
	if err != nil {
		r.Error = err
		b.report = r
		return
	}

	ctx := &Context{
		DstDir:          dest,
		Theme:           theme,
		Version:         b.Version,
		isCopyAllAssets: b.opt.IsCopyAssets,
		isSuffixed:      b.opt.IsSuffixed,
	}

	b.isBuilding = true

	// run tasks
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
		b.report = r
		b.context = ctx
	}

	log15.Info("Build.Finish", "duration", r.Duration(), "error", r.Error)
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

// get last report in builder
func (b *Builder) Report() *Report {
	return b.report
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
