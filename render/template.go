package render

import (
	"errors"
	"html/template"
	"io"
	"path"
)

var (
	ErrTemplateMissing = errors.New("template-missing")
)

type Template struct {
	file   string
	tpl    *template.Template
	reload bool
	Error  error
}

func newTemplate(file string, reload bool) *Template {
	return &Template{
		file:   file,
		reload: reload,
	}
}

func (t *Template) Compile(w io.Writer, data interface{}, fnMap template.FuncMap) {
	if t.tpl == nil || t.reload {
		tpl, err := template.New(path.Base(t.file)).Funcs(fnMap).ParseFiles(t.file)
		if err != nil {
			t.Error = err
			return
		}
		t.tpl = tpl
	}
	if err := t.tpl.Execute(w, data); err != nil {
		t.Error = err
		return
	}
}
