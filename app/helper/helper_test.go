package helper_test

import (
	"html/template"
	"testing"

	"github.com/go-xiaohei/pugo/app/helper"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPager(t *testing.T) {
	cursor := helper.NewPagerCursor(5, 99)
	Convey("test pager", t, func() {
		Convey("test all pager", func() {
			var i = 1
			for i = 1; i < 100; i++ {
				if p := cursor.Page(i); p == nil {
					break
				}
			}
			So(i-1, ShouldEqual, 20)

			p := cursor.Page(-1)
			So(p, ShouldBeNil)
		})

		Convey("each pager", func() {
			p := cursor.Page(1)
			So(p.Prev, ShouldEqual, 0)
			So(p.PrevURL(), ShouldBeEmpty)

			p = cursor.Page(20)
			So(p.Prev, ShouldEqual, 19)
			So(p.Next, ShouldEqual, 0)
			So(p.End, ShouldEqual, 99)
			So(p.NextURL(), ShouldBeEmpty)

			p = cursor.Page(10)
			p.SetLayout("page%d")
			So(p.PrevURL(), ShouldEqual, "page9")
			So(p.NextURL(), ShouldEqual, "page11")
		})
	})
}

func TestI18n(t *testing.T) {
	i18n, err := helper.NewI18n("../../source/lang/en.ini", "")
	Convey("test i18n init", t, func() {
		So(err, ShouldBeNil)
	})

	Convey("test i18n translate", t, func() {
		v := i18n.Tr("post.list")
		So(v, ShouldEqual, "All Posts")

		v = i18n.Trf("post.list")
		So(v, ShouldEqual, "All Posts")

		v2 := i18n.TrHTML("post.list")
		So(v2, ShouldHaveSameTypeAs, template.HTML(""))

		v2 = i18n.TrfHTML("post.list")
		So(v2, ShouldHaveSameTypeAs, template.HTML(""))

		v3 := i18n.Tr("post.null")
		So(v3, ShouldEqual, "post.null")
	})

	Convey("test i18n inti error", t, func() {
		_, err := helper.NewI18n("xxx.ini", "")
		So(err, ShouldNotBeNil)
	})

}
