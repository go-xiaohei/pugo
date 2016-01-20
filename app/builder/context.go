package builder

import "github.com/go-xiaohei/pugo/app/theme"

type (
	// Context obtain context in once building process
	Context struct {
		// From is source origin
		From string
		// To is destination
		To string
		// Theme is theme origin
		ThemeName string
		// Err is error when context using
		Err error
		// Source is sources data
		Source *Source
		// Theme is theme object, use to render templates
		Theme *theme.Theme
	}
)

// NewContext create new Context with from,to and theme args
func NewContext(from, to, theme string) *Context {
	return &Context{
		From:      from,
		To:        to,
		ThemeName: theme,
	}
}

// IsValid check context requirement, there must have values in some fields
func (ctx *Context) IsValid() bool {
	if ctx.From == "" || ctx.To == "" || ctx.ThemeName == "" {
		return false
	}
	return true
}
