package render_test

import (
	"bytes"
	"github.com/go-xiaohei/pugo-static/render"
	. "github.com/smartystreets/goconvey/convey"
	"html/template"
	"testing"
)

func TestRender(t *testing.T) {
	r := render.New("../template")
	Convey("render load theme", t, func() {
		theme, err := r.Load("default")
		So(err, ShouldBeNil)
		So(theme, ShouldHaveSameTypeAs, new(render.Theme))

		_, err = r.Load("xxxxxx")
		So(err, ShouldEqual, render.ErrRenderDirMissing)
	})
}

func TestTheme(t *testing.T) {
	funcMap := make(template.FuncMap)
	funcMap["include"] = func(tpl string) template.HTML {
		return template.HTML(tpl)
	}
	theme := render.NewTheme("../template/default", funcMap, []string{".html"})
	Convey("load theme", t, func() {
		err := theme.Load()
		So(err, ShouldBeNil)

		Convey("render file", func() {
			var buf bytes.Buffer
			err := theme.Execute(&buf, "xxxxxx.html", nil)
			So(err, ShouldEqual, render.ErrTemplateMissing)

			buf.Reset()
			err = theme.Execute(&buf, "post.html", nil)
			So(err, ShouldBeNil)

			buf.Reset()
			err = theme.Execute(&buf, "posts.html", nil)
			So(err, ShouldBeNil)
		})
	})

}
