package model

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFile(t *testing.T) {
	files := NewFiles()

	Convey("Files", t, func() {
		files.Add("abc.html", 123, time.Now(), FileCompiled, OpCompiled)
		files.Add("xyz.html", 456, time.Now(), FileStatic, OpRemove)

		So(files.Get("xyz.html"), ShouldNotBeNil)
		So(files.Exist("xyz.jpg"), ShouldBeFalse)
		So(files.All(), ShouldHaveLength, 2)
	})
}
