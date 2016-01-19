package builder

type (
	// Context obtain context in once building process
	Context struct {
		// From is source origin
		From string
		// To is destination
		To string
		// Theme is theme origin
		Theme string
		// Err is error when context using
		Err error
		// Source is sources data
		Source *Source
	}
)

// NewContext create new Context with from,to and theme args
func NewContext(from, to, theme string) *Context {
	return &Context{
		From:  from,
		To:    to,
		Theme: theme,
	}
}

// IsValid check context requirement, there must have values in some fields
func (ctx *Context) IsValid() bool {
	if ctx.From == "" || ctx.To == "" || ctx.Theme == "" {
		return false
	}
	return true
}
