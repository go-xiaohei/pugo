package builder

import (
	"errors"
	"github.com/Unknwon/com"
	"pugo/render"
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

func (b *Builder) Build() {
	// if on build, do not again
	if b.isBuilding {
		return
	}
	b.isBuilding = true
}

func (b *Builder) IsBuilding() bool {
	return b.isBuilding
}

func (b *Builder) Report() *Report {
	return b.report
}
