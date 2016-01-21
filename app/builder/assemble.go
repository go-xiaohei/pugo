package builder

import (
	"fmt"
	"github.com/go-xiaohei/pugo/app/model"
	"path"
	"strings"
)

// AssembleSource assemble some extra data in Source,
// such as page nodes, i18n status.
// it need be used after posts and pages are loaded
func AssembleSource(ctx *Context) {
	if ctx.Source == nil || ctx.Theme == nil {
		ctx.Err = fmt.Errorf("need sources data and theme to assemble")
		return
	}

	ctx.Source.Nav.FixURL(ctx.Source.Meta.Path)
	ctx.Source.Tree = model.NewTree()

	r, hr := newReplacer("/"+ctx.Theme.Static()), newReplacerInHTML("/"+ctx.Theme.Static())
	if ctx.Source.Meta.Path != "" && ctx.Source.Meta.Path != "/" {
		for _, p := range ctx.Source.Posts {
			p.FixURL(ctx.Source.Meta.Path)
			p.FixPlaceholder(r, hr)
			p.Author = ctx.Source.Authors[p.AuthorName]
			ctx.Source.Tree.Add(p.TreeURL(), model.TreePost, 0)
		}
		for _, p := range ctx.Source.Pages {
			p.FixURL(ctx.Source.Meta.Path)
			p.FixPlaceholder(hr)
			p.Author = ctx.Source.Authors[p.AuthorName]
			ctx.Source.Tree.Add(p.TreeURL(), model.TreePage, p.Sort)
		}
	}
	if ctx.Err = ctx.Theme.Load(); ctx.Err != nil {
		return
	}
}

func newReplacer(static string) *strings.Replacer {
	return strings.NewReplacer(
		"@static", static,
		"@media", path.Join(static, "media"),
	)
}

func newReplacerInHTML(static string) *strings.Replacer {
	media := path.Join(static, "media")
	return strings.NewReplacer(
		`src="@static`, fmt.Sprintf(`src="%s`, static),
		`href="@static`, fmt.Sprintf(`href="%s`, static),
		`src="@media`, fmt.Sprintf(`src="%s`, media),
		`href="@media`, fmt.Sprintf(`src="%s`, media),
	)
}
