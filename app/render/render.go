package render

import (
	"fmt"
	"html/template"
	"path"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/helper"
)

const (
	ErrRenderDirMissing int = 1 // error of template directory missing
	ErrTemplateMissing  int = 2 // error of template file missing
)

var (
	errorDirMissing = RenderError{
		Message: "template dir '%s' is missing",
		Type:    ErrRenderDirMissing,
	}
	errorFileMissing = RenderError{
		Message: "template file '%s' is missing",
		Type:    ErrTemplateMissing,
	}
)

// render struct
type (
	Render struct {
		dir        string
		extensions []string
		funcMap    template.FuncMap
	}
	RenderError struct {
		Message string
		Type    int
	}
)

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
		return nil, errorDirMissing.New(name)
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

// error implementation
func (err RenderError) Error() string {
	return err.Message
}

// new message error
func (err RenderError) New(values ...interface{}) RenderError {
	err.Message = fmt.Sprintf(err.Message, values...)
	return err
}
