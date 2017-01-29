package helper

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMd5(t *testing.T) {
	Convey("Md5", t, func() {
		str := Md5("123456")
		So(str, ShouldEqual, "e10adc3949ba59abbe56e057f20f883e")

		str, err := Md5File("md5.go")
		So(err, ShouldBeNil)
		So(str, ShouldEqual, "651e74ed7f68be2b642217a06fda6ec6")
	})
}
