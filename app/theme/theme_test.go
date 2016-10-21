package theme

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTheme(t *testing.T) {
	theme := New("../../theme/default")

	Convey("LoadTheme", t, func() {
		err := theme.Load()
		So(err, ShouldBeNil)

		Convey("Funcs", func() {
			funcs := theme.Funcs()
			So(funcs, ShouldContainKey, "HTML")
			So(funcs, ShouldContainKey, "Include")
		})
	})

}

func TestThemeMeta(t *testing.T) {
	theme := New("../../theme/uno")

	Convey("LoadTheme", t, func() {
		err := theme.Load()
		So(err, ShouldBeNil)
		So(theme.Validate(), ShouldBeNil)

		Convey("Authors", func() {
			So(theme.Meta.Authors, ShouldHaveLength, 1)
			So(theme.Meta.Authors[0].Name, ShouldEqual, "fuxiaohei")
		})

		Convey("Ref", func() {
			So(theme.Meta.Refs, ShouldHaveLength, 1)
			So(theme.Meta.Refs[0].Name, ShouldEqual, "hexo uno theme")
		})
	})

}
