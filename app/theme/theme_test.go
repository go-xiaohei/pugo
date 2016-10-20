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
