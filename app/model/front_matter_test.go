package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFrontMatter(t *testing.T) {
	Convey("FrontMetaFileNoExist", t, func() {
		_, err := NewPostsFrontMatter("testdata/metaaaaa.toml", FormatTOML)
		So(err, ShouldNotBeNil)

		_, err = NewPostsFrontMatter("testdata/meta_wrong.toml", FormatTOML)
		So(err, ShouldNotBeNil)

		_, err = NewPagesFrontMatter("testdata/metaaaaa.ini", FormatTOML)
		So(err, ShouldNotBeNil)

		_, err = NewPagesFrontMatter("testdata/meta_wrong.ini", FormatTOML)
		So(err, ShouldNotBeNil)
	})
}
