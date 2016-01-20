package builder

import (
	"fmt"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app2/theme"

	"gopkg.in/inconshreveable/log15.v2"
)

// ReadTheme read *Theme to *Context
func ReadTheme(ctx *Context) {
	// TODO: clean code
	ctx.ThemeName = "../theme/default"
	if ctx.Source == nil {
		ctx.Err = fmt.Errorf("theme depends on loaded source data")
		return
	}
	if !com.IsDir(ctx.ThemeName) {
		ctx.Err = fmt.Errorf("theme directory '%s' is missing", ctx.ThemeName)
		return
	}
	log15.Debug("Build|Theme|%s", ctx.ThemeName)
	ctx.Theme = theme.New(ctx.ThemeName)
}
