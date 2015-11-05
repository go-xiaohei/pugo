package render

import "errors"

var (
	ErrTemplateMissing = errors.New("template-missing")
)

type Template struct {
	Error error
}

func newTemplate(file string) *Template {
	return nil
}
