package builder

import (
	"fmt"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/theme"
	"gopkg.in/inconshreveable/log15.v2"
)

// ReadTheme read *Theme to *Context
func ReadTheme(ctx *Context) {
	if ctx.Source == nil {
		ctx.Err = fmt.Errorf("theme depends on loaded source data")
		return
	}
	dir, _ := toDir(ctx.ThemeName)
	if !com.IsDir(dir) {
		ctx.Err = fmt.Errorf("theme directory '%s' is missing", dir)
		return
	}
	log15.Info("Build|Theme|%s", dir)
	ctx.Theme = theme.New(dir)
}
