package render

import (
	"errors"
	"html/template"
	"path"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/helper"
)

var (
	ErrRenderDirMissing = errors.New("render-dir-missing") // error of template directory missing
)

// render struct
type Render struct {
	dir        string
	extensions []string
	funcMap    template.FuncMap
}

// new render in directory
func New(dir string) *Render {
	r := &Render{
		dir:        dir,
		extensions: []string{".html"},
		funcMap:    make(template.FuncMap),
	}
	r.funcMap["HTML"] = helper.Str2HTML
	r.funcMap["HTMLbyte"] = helper.Bytes2HTML
	return r
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
