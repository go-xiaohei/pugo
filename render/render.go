package render

import (
	"bytes"
	"fmt"
	"github.com/Unknwon/com"
	"html/template"
	"path"
)

type (
	Render struct {
		dir       string
		templates map[string]*Template
		reload    bool
		fnMap     template.FuncMap
	}
)

func NewRender(dir string, reload bool) *Render {
	return &Render{
		dir:       dir,
		templates: make(map[string]*Template),
		reload:    reload,
	}
}

func (r *Render) StaticDir() string {
	return path.Join(r.dir, "static")
}

func (r *Render) Template(file string) *Template {
	if r.templates[file] == nil {
		fullFile := path.Join(r.dir, file)
		if !com.IsFile(fullFile) {
			return &Template{Error: ErrTemplateMissing}
		}
		template := newTemplate(fullFile, r.reload)
		if template.Error != nil {
			return template
		}
		r.templates[file] = template
	}
	return r.templates[file]
}

func (r *Render) FuncMap() template.FuncMap {
	if r.fnMap == nil {
		fnMap := make(template.FuncMap)
		fnMap["include"] = func(file string, data interface{}) template.HTML {
			t := r.Template(file)
			if t.Error != nil {
				return template.HTML(fmt.Sprintf("<-- include %s error : %s -->", file, t.Error.Error()))
			}
			var buf bytes.Buffer
			if t.Compile(&buf, data, r.fnMap); t.Error != nil {
				return template.HTML(fmt.Sprintf("<-- include %s error : %s -->", file, t.Error.Error()))
			}
			return template.HTML(buf.Bytes())
		}
		r.fnMap = fnMap
	}
	return r.fnMap
}
