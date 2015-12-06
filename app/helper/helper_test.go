package helper_test

import (
	"testing"

	"github.com/go-xiaohei/pugo-static/app/helper"
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
			So(p.PrevUrl(), ShouldBeEmpty)

			p = cursor.Page(20)
			So(p.Prev, ShouldEqual, 19)
			So(p.Next, ShouldEqual, 0)
			So(p.End, ShouldEqual, 99)
			So(p.NextUrl(), ShouldBeEmpty)

			p = cursor.Page(10)
			p.SetLayout("page%d")
			So(p.PrevUrl(), ShouldEqual, "page9")
			So(p.NextUrl(), ShouldEqual, "page11")
		})
	})
}
