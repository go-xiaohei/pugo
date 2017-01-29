package helper

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPager(t *testing.T) {
	Convey("Pager", t, func() {
		pager := NewPagerCursor(5, 99)

		Convey("Page", func() {
			page := pager.Page(3)
			So(page.Begin, ShouldEqual, 10)
			So(page.End, ShouldEqual, 15)

			page = pager.Page(20)
			So(page.End, ShouldEqual, 99)

			page = pager.Page(-1)
			So(page, ShouldBeNil)

			page = pager.Page(1000)
			So(page, ShouldBeNil)
		})

		Convey("Layout", func() {
			page := pager.Page(3)
			page.SetLayout("ppp%d")
			So(page.URL(), ShouldEqual, "ppp3")

			So(page.PrevURL(), ShouldEqual, "ppp2")
			So(page.NextURL(), ShouldEqual, "ppp4")

			page = pager.Page(1)
			So(page.PrevURL(), ShouldEqual, "")

			page = pager.Page(20)
			So(page.NextURL(), ShouldEqual, "")
		})

		Convey("SizeCase", func() {
			pager := NewPagerCursor(10, 100)
			So(pager.pages, ShouldEqual, 10)
			pager = NewPagerCursor(10, 98)
			So(pager.pages, ShouldEqual, 10)
		})

		Convey("PagerItems", func() {
			pager := NewPagerCursor(10, 100)
			page := pager.Page(6)
			page.SetLayout("aaa%d")
			for i, item := range page.PageItems() {
				So(item.Page, ShouldEqual, i+1)
				So(item.Link, ShouldEqual, fmt.Sprintf("aaa%d", i+1))
			}
		})
	})
}
