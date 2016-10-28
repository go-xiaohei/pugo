package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFormat(t *testing.T) {
	Convey("TestFormatExtension", t, func() {
		suffixData := ShouldFormatExtension()
		So(suffixData, ShouldHaveLength, 2)
		for t, suffix := range suffixData {
			if t == FormatTOML {
				So(suffix, ShouldEqual, ".toml")
			}
			if t == FormatINI {
				So(suffix, ShouldEqual, ".ini")
			}
		}
	})
	Convey("TestThemeMetaFiles", t, func() {
		suffixData := ShouldThemeMetaFiles()
		So(suffixData, ShouldHaveLength, 2)
		for t, suffix := range suffixData {
			if t == FormatTOML {
				So(suffix, ShouldEqual, "theme.toml")
			}
			if t == FormatINI {
				So(suffix, ShouldEqual, "theme.ini")
			}
		}
	})
}
