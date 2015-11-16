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
	Builder struct {
		srcDir     string
		tplDir     string
		isBuilding bool
		opt        *BuildOption

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

func New(opt *BuildOption) *Builder {
	if !com.IsDir(opt.SrcDir) {
		return &Builder{Error: ErrSrcDirMissing}
	}
	if !com.IsDir(opt.TplDir) {
		return &Builder{Error: ErrTplDirMissing}
	}
	builder := &Builder{
		srcDir: opt.SrcDir,
		tplDir: opt.TplDir,
		parser: parser.NewCommonParser(),
		Version: builderVersion{
			Num:  opt.Version,
			Date: opt.VerDate,
		},
		opt: opt,
	}
	r, err := render.NewRenders(opt.TplDir, opt.Theme, opt.IsDebug)
	if err != nil {
		return &Builder{Error: err}
	}
	builder.renders = r
	builder.tasks = []*BuildTask{
		&BuildTask{"Data", builder.ReadData, nil},
		&BuildTask{"Compile", builder.Compile, nil},
		&BuildTask{"Copy", builder.CopyAssets, nil},
		&BuildTask{"Feed", builder.WriteFeed, nil},
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
		DstDir:          dest,
		Version:         b.Version,
		isCopyAllAssets: b.opt.IsCopyAssets,
		isSuffixed:      b.opt.IsSuffixed,
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
