package render

import (
	"github.com/Unknwon/com"
	"path"
)

type (
	Render struct {
		dir       string
		templates map[string]*Template
		reload    bool
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
		template := newTemplate(fullFile)
		if template.Error != nil {
			return template
		}
		r.templates[file] = template
	}
	return r.templates[file]
}
