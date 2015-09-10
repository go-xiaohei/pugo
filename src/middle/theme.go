package middle

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
	"path"
	"pugo/src/core"
	"pugo/src/model"
	"pugo/src/service"
	"strings"
)

var (
	ThemeErrorTemplate = "error.tmpl"

	_ ITheme = (*baseTheme)(nil)
)

type ITheme interface {
	SetTheme(*model.Theme)
	GetTheme() *model.Theme

	Assign(key string, value interface{})
	Title(string)
	IsAssign(key string) bool
	Render(tpl string)
	RenderError(status int, err error)
}

type baseTheme struct {
	renders.Renderer

	themeDirectory string
	data           map[string]interface{}
	theme          *model.Theme
}

func (bt *baseTheme) SetTheme(theme *model.Theme) {
	bt.theme = theme
	bt.themeDirectory = theme.Directory
}

func (bt *baseTheme) GetTheme() *model.Theme {
	return bt.theme
}

func (bt *baseTheme) Assign(key string, value interface{}) {
	if len(bt.data) == 0 {
		bt.data = make(map[string]interface{})
		bt.data["ThemeLink"] = path.Join(core.ThemePrefix, strings.TrimPrefix(bt.themeDirectory, core.ThemeDirectory))
	}
	bt.data[key] = value
}

func (bt *baseTheme) Title(title string) {
	bt.Assign("Title", title)
}

func (bt *baseTheme) IsAssign(key string) bool {
	_, ok := bt.data[key]
	return ok
}

func (bt *baseTheme) Render(tpl string) {
	tpl = path.Join(bt.themeDirectory, tpl)
	if err := bt.Renderer.Render(tpl, bt.data); err != nil {
		panic(err)
	}
}

func (bt *baseTheme) RenderError(status int, err error) {
	tpl := path.Join(bt.themeDirectory, ThemeErrorTemplate)
	bt.Assign("Error", err)
	if err := bt.StatusRender(status, tpl, bt.data); err != nil {
		panic(err)
	}
}

type ThemeRender struct {
	baseTheme
}

func (t *ThemeRender) SetTheme(*model.Theme) {
	var theme = new(model.Theme)
	if err := service.Call(service.Theme.Current, nil, theme); err != nil {
		panic(err)
	}
	t.baseTheme.SetTheme(theme)
}

type AdminRender struct {
	baseTheme
}

func (t *AdminRender) SetTheme(*model.Theme) {
	var theme = new(model.Theme)
	if err := service.Call(service.Theme.Admin, nil, theme); err != nil {
		panic(err)
	}
	t.baseTheme.SetTheme(theme)
}

func Themer() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		t, ok := ctx.Action().(ITheme)
		if ok {
			t.SetTheme(nil) // set theme, the implementions finish the method in local scope
		}
		ctx.Next()
	}
}
