package middle

import (
	"fmt"
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
	"path"
	"pugo/src/core"
	"pugo/src/model"
	"pugo/src/service"
	"reflect"
	"strings"
)

var (
	ThemeErrorTemplate = "error.tmpl"

	_ ITheme = (*baseTheme)(nil)
	_ ITheme = (*ThemeRender)(nil)
	_ ITheme = (*AdminRender)(nil)
)

type ITheme interface {
	SetTheme(*model.Theme) error
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

func (bt *baseTheme) SetTheme(theme *model.Theme) error {
	bt.theme = theme
	bt.themeDirectory = theme.Directory
	return nil
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
	tpl = strings.TrimPrefix(tpl, core.ThemeDirectory+"/")
	if err := bt.Renderer.Render(tpl, bt.data); err != nil {
		panic(err)
	}
}

func (bt *baseTheme) RenderError(status int, err error) {
	tpl := path.Join(bt.themeDirectory, ThemeErrorTemplate)
	tpl = strings.TrimPrefix(tpl, core.ThemeDirectory+"/")
	bt.Assign("Error", err)
	if err := bt.StatusRender(status, tpl, bt.data); err != nil {
		panic(err)
	}
}

type ThemeRender struct {
	baseTheme
}

func (t *ThemeRender) SetTheme(*model.Theme) error {
	var theme = new(model.Theme)
	if err := service.Call(service.Theme.Current, nil, theme); err != nil {
		return err
	}
	t.baseTheme.SetTheme(theme)
	return nil
}

type AdminRender struct {
	baseTheme
}

func (t *AdminRender) SetTheme(*model.Theme) error {
	var theme = new(model.Theme)
	if err := service.Call(service.Theme.Admin, nil, theme); err != nil {
		return err
	}
	t.baseTheme.SetTheme(theme)
	return nil
}

func Themer() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		t, ok := ctx.Action().(ITheme)
		if ok {
			if err := t.SetTheme(nil); err != nil { // set theme, the implement finish the method in local scope
				ctx.Result = fmt.Errorf("%s %s", reflect.TypeOf(t).String(), err.Error())
				ctx.HandleError()
				return
			}
		}
		ctx.Next()
	}
}
