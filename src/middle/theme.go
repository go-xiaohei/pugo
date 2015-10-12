package middle

import (
	"fmt"
	"github.com/go-xiaohei/pugo/src/core"
	"github.com/go-xiaohei/pugo/src/model"
	"github.com/go-xiaohei/pugo/src/service"
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
	"net/http"
	"path"
	"reflect"
	"strings"
)

var (
	ThemeErrorTemplate = "error.tmpl"

	_ ITheme = (*baseTheme)(nil)
	_ ITheme = (*TemplateRender)(nil)
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
	if theme == nil {
		bt.theme = nil
		bt.themeDirectory = ""
		return nil
	}
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
	bt.Assign("Title", fmt.Sprintf("%d - %s", status, strings.ToUpper(http.StatusText(status))))
	bt.Assign("Status", status)
	bt.Assign("Error", err)
	if err := bt.StatusRender(status, tpl, bt.data); err != nil {
		panic(err)
	}
}

type TemplateRender struct {
	baseTheme
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

func (t *ThemeRender) Render(tpl string) {
	if !t.IsAssign("General") {
		t.Assign("General", service.Setting.General)
		t.Assign("Menu", service.Setting.Menu)
	}
	t.baseTheme.Render(tpl)
}

func (t *ThemeRender) RenderError(status int, err error) {
	if !t.IsAssign("General") {
		t.Assign("General", service.Setting.General)
		t.Assign("Menu", service.Setting.Menu)
	}
	t.baseTheme.RenderError(status, err)
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
