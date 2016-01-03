package render

import (
	"fmt"
	"html/template"
	"path"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/helper"
)

const (
	// ErrRenderDirMissing is error of template directory missing
	ErrRenderDirMissing int = 1
	// ErrTemplateMissing is error of template file missing
	ErrTemplateMissing int = 2
)

var (
	errorDirMissing = Error{
		Message: "template dir '%s' is missing",
		Type:    ErrRenderDirMissing,
	}
	errorFileMissing = Error{
		Message: "template file '%s' is missing",
		Type:    ErrTemplateMissing,
	}
)

type (
	// Render struct
	Render struct {
		dir        string
		extensions []string
		funcMap    template.FuncMap
	}
	// Error is error type of render
	Error struct {
		Message string
		Type    int
	}
)

// New returns new render in directory
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

// Load loads a theme by name
func (r *Render) Load(name string) (*Theme, error) {
	dir := path.Join(r.dir, name)
	if !com.IsDir(dir) {
		return nil, errorDirMissing.New(name)
	}
	theme := NewTheme(dir, r.funcMap, r.extensions)
	return theme, theme.Load()
}

// SetExtensions sets extensions
func (r *Render) SetExtensions(ex []string) {
	r.extensions = ex
}

// SetFunc sets func map by name
func (r *Render) SetFunc(name string, fn interface{}) {
	if fn == nil {
		delete(r.funcMap, name)
		return
	}
	r.funcMap[name] = fn
}

// Error is error implementation
func (err Error) Error() string {
	return err.Message
}

// New returns new message error
func (err Error) New(values ...interface{}) Error {
	err.Message = fmt.Sprintf(err.Message, values...)
	return err
}
