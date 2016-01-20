package builder

import "fmt"

// AssembleSource assemble some extra data in Source,
// such as page nodes, i18n status.
// it need be used after posts and pages are loaded
func AssembleSource(ctx *Context) {
	if ctx.Source == nil || ctx.Theme == nil {
		ctx.Err = fmt.Errorf("need sources data and theme to assemble")
		return
	}
	if ctx.Err = ctx.Theme.Load(); ctx.Err != nil {
		return
	}
}
