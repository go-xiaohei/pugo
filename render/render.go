package render

import (
	"errors"
	"github.com/Unknwon/com"
	"html/template"
	"path"
)

var (
	ErrRenderDirMissing = errors.New("render-dir-missing")
)

type Render struct {
	dir        string
	extensions []string
	funcMap    template.FuncMap
}

func New(dir string) *Render {
	return &Render{
		dir:        dir,
		extensions: []string{".html"},
		funcMap:    make(template.FuncMap),
	}
}

// load theme by name
func (r *Render) Load(name string) (*Theme, error) {
	dir := path.Join(r.dir, name)
	if !com.IsDir(dir) {
		return nil, ErrRenderDirMissing
	}
	theme := NewTheme(dir, r.funcMap, r.extensions)
	return theme, theme.Load()
}

// set extensions
func (r *Render) SetExtensions(ex []string) {
	r.extensions = ex
}

// set func map by name
func (r *Render) SetFunc(name string, fn interface{}) {
	if fn == nil {
		delete(r.funcMap, name)
		return
	}
	r.funcMap[name] = fn
}
