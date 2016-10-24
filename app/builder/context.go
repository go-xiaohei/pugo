package builder

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/go-xiaohei/pugo/app/model"
	"github.com/go-xiaohei/pugo/app/theme"
	"github.com/go-xiaohei/pugo/app/vars"
	"github.com/urfave/cli"
	"gopkg.in/inconshreveable/log15.v2"
)

type (
	// Context obtain context in once building process
	Context struct {
		cli *cli.Context
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
func NewContext(cli *cli.Context, from, to, theme string) *Context {
	c := &Context{
		cli:       cli,
		From:      from,
		To:        to,
		ThemeName: theme,
		time:      time.Now(),
		Files:     model.NewFiles(),
		Copied:    defaultCopiedOpt(),
	}
	c.Tree = model.NewTree(c.DstDir())
	return c
}

// View get view data to template from Context
func (ctx *Context) View() map[string]interface{} {
	m := map[string]interface{}{
		"Version":   vars.Version,
		"Source":    ctx.Source,
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
		m["I18n"] = ctx.Source.I18n["en"]
		if m["I18n"] == nil {
			m["I18n"] = helper.NewI18nEmpty()
		}
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

// SrcDir get src dir
func (ctx *Context) SrcDir() string {
	ctx.parseDir()
	return ctx.srcDir
}

// SrcPostDir get post dir in src
func (ctx *Context) SrcPostDir() string {
	ctx.parseDir()
	if ctx.Source != nil && ctx.Source.Build != nil && ctx.Source.Build.PostDir != "" {
		return path.Join(ctx.srcDir, ctx.Source.Build.PostDir)
	}
	return path.Join(ctx.srcDir, "post")
}

// SrcPageDir get page dir in src
func (ctx *Context) SrcPageDir() string {
	ctx.parseDir()
	if ctx.Source != nil && ctx.Source.Build != nil && ctx.Source.Build.PageDir != "" {
		return path.Join(ctx.srcDir, ctx.Source.Build.PageDir)
	}
	return path.Join(ctx.srcDir, "page")
}

// SrcLangDir get language dir in src
func (ctx *Context) SrcLangDir() string {
	ctx.parseDir()
	if ctx.Source != nil && ctx.Source.Build != nil && ctx.Source.Build.LangDir != "" {
		return path.Join(ctx.srcDir, ctx.Source.Build.LangDir)
	}
	return path.Join(ctx.srcDir, "lang")
}

// SrcThemeDir get theme dir in src
func (ctx *Context) SrcThemeDir() string {
	ctx.parseDir()
	if ctx.ThemeName != "" {
		return ctx.ThemeName
	}
	if ctx.Source != nil && ctx.Source.Build != nil && ctx.Source.Build.ThemeDir != "" {
		return ctx.Source.Build.ThemeDir
	}
	return "source/theme/default"
}

// DstDir get destination directory after build once
func (ctx *Context) DstDir() string {
	ctx.parseDir()
	return ctx.dstDir
}

// Cli get command line context in this building context
func (ctx *Context) Cli() *cli.Context {
	return ctx.cli
}

func (ctx *Context) parseDir() {
	if ctx.srcDir != "" && ctx.dstDir != "" {
		return
	}
	var (
		srcDir  = ""
		destDir = ""
	)
	if srcDir, ctx.Err = toDir(ctx.From); ctx.Err != nil {
		return
	}
	if !com.IsDir(srcDir) {
		ctx.Err = fmt.Errorf("Directory '%s' is missing", srcDir)
		return
	}
	ctx.srcDir = srcDir
	log15.Info("Read|%s", srcDir)

	if destDir, ctx.Err = toDir(ctx.To); ctx.Err != nil {
		return
	}
	ctx.dstDir = destDir
}

func toDir(urlString string) (string, error) {
	if !strings.Contains(urlString, "://") {
		return urlString, nil
	}
	if strings.HasPrefix(urlString, "dir://") {
		return strings.TrimPrefix(urlString, "dir://"), nil
	}
	return "", errors.New("Directory need schema dir://")
}
