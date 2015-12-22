package render_test

import (
	"bytes"
	"html/template"
	"path"
	"testing"

	"github.com/go-xiaohei/pugo/app/model"
	"github.com/go-xiaohei/pugo/app/render"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	funcMap template.FuncMap
)

func init() {
	funcMap = make(template.FuncMap)
	funcMap["include"] = func(tpl string) template.HTML {
		return template.HTML(tpl)
	}
	funcMap["url"] = func(str ...string) string {
		return path.Join(str...)
	}
	funcMap["fullUrl"] = funcMap["url"]
}

func TestRender(t *testing.T) {
	r := render.New("../../template")
	for name, fn := range funcMap {
		r.SetFunc(name, fn)
	}
	Convey("render load theme", t, func() {
		theme, err := r.Load("default")
		So(err, ShouldBeNil)
		So(theme, ShouldHaveSameTypeAs, new(render.Theme))

		_, err = r.Load("xxxxxx")
		So(err, ShouldEqual, render.ErrRenderDirMissing)
	})
}

func TestTheme(t *testing.T) {
	theme := render.NewTheme("../../template/default", funcMap, []string{".html"})
	Convey("load theme", t, func() {
		err := theme.Load()
		So(err, ShouldBeNil)

		Convey("render file", func() {
			var buf bytes.Buffer
			err := theme.Execute(&buf, "xxxxxx.html", nil)
			So(err, ShouldEqual, render.ErrTemplateMissing)

			buf.Reset()
			err = theme.Execute(&buf, "post.html", map[string]interface{}{
				"Post": &model.Post{Author: new(model.Author)},
			})
			So(err, ShouldBeNil)

			buf.Reset()
			err = theme.Execute(&buf, "posts.html", map[string]interface{}{
				"Posts": []*model.Post{},
			})
			So(err, ShouldBeNil)
		})
	})

}
