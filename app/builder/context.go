package builder

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/go-xiaohei/pugo/app/model"
	"github.com/go-xiaohei/pugo/app/theme"
	"github.com/go-xiaohei/pugo/app/vars"
)

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
		// Files is generated files in by this context
		Files *model.Files
		// Tree is url tree nodes by this context
		Tree   *model.Tree
		Copied *CopiedOpt

		time           time.Time
		counter        int64
		srcDir, dstDir string
	}
)

// NewContext create new Context with from,to and theme args
func NewContext(from, to, theme string) *Context {
	return &Context{
		From:      from,
		To:        to,
		ThemeName: theme,
		time:      time.Now(),
		Files:     model.NewFiles(),
		Copied:    defaultCopiedOpt(),
		Tree:      model.NewTree(),
	}
}

// View get view data to template from Context
func (ctx *Context) View() map[string]interface{} {
	m := map[string]interface{}{
		"Version":   vars.Version,
		"Nav":       ctx.Source.Nav,
		"Meta":      ctx.Source.Meta,
		"Title":     ctx.Source.Meta.Title + " - " + ctx.Source.Meta.Subtitle,
		"Desc":      ctx.Source.Meta.Desc,
		"Comment":   ctx.Source.Comment,
		"Owner":     ctx.Source.Owner,
		"Analytics": ctx.Source.Analytics,
		"Tree":      ctx.Tree,
		"Lang":      ctx.Source.Meta.Language,
		"Hover":     "",
		"Base":      strings.TrimRight(ctx.Source.Meta.Path, "/"),
		"Root":      strings.TrimRight(ctx.Source.Meta.Root, "/"),
	}
	if ctx.Source.Meta.Language == "" {
		m["I18n"] = helper.NewI18nEmpty()
	} else {
		if i18n, ok := ctx.Source.I18n[ctx.Source.Meta.Language]; ok {
			m["I18n"] = i18n
		} else {
			m["I18n"] = helper.NewI18nEmpty()
		}
	}
	return m
}

// IsValid check context requirement, there must have values in some fields
func (ctx *Context) IsValid() bool {
	if ctx.From == "" || ctx.To == "" || ctx.ThemeName == "" {
		return false
	}
	return true
}

// Duration return seconds after *Context created
func (ctx *Context) Duration() float64 {
	return time.Since(ctx.time).Seconds()
}

// Again reset some fields in context to rebuild
func (ctx *Context) Again() {
	ctx.time = time.Now()
	atomic.StoreInt64(&ctx.counter, 0)
}

// SrcDir get src dir after build once
func (ctx *Context) SrcDir() string {
	return ctx.srcDir
}

// DstDir get destination directory after build once
func (ctx *Context) DstDir() string {
	return ctx.dstDir
}
