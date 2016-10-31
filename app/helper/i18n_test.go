package helper

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	i18nBytes = []byte(`[meta]
title = "Title"
subtitle = "Subtitle"
link = "Link%s"

[nav]
home = "Home"
archive = "Archive"
about = "About"
source = "Source"`)

	i18nIniBytes = []byte(`
"meta.link" = "Link%s"

[meta]
title = "Title"
subtitle = "Subtitle"

[nav]
home = "Home"
archive = "Archive"
about = "About"
source = "Source"
`)
)

func TestI18n(t *testing.T) {
	Convey("I18n", t, func() {
		i18n, err := NewI18n("en", i18nBytes, ".toml")
		So(err, ShouldBeNil)

		Convey("Tr", func() {
			tr := i18n.Tr("meta.title")
			So(tr, ShouldEqual, "Title")

			tr = i18n.Tr("meta.xxx")
			So(tr, ShouldEqual, "meta.xxx")

			tr = i18n.Tr("a.b.c")
			So(tr, ShouldEqual, "a.b.c")
		})

		Convey("Trf", func() {
			tr := i18n.Trf("meta.link", "abc")
			So(tr, ShouldEqual, "Linkabc")
		})

		Convey("UnmashalFail", func() {
			b := []byte(`abc="abc"`)
			_, err := NewI18n("en", b, ".toml")
			So(err, ShouldNotBeNil)
		})

		Convey("Empty", func() {
			i18n := NewI18nEmpty()
			So(i18n.values, ShouldHaveLength, 0)

			tr := i18n.Tr("meta.title")
			So(tr, ShouldEqual, "meta.title")
		})

		Convey("Lang", func() {
			en := "en-US"
			codes := LangCode(en)
			So(codes, ShouldHaveLength, 3)
			So(codes, ShouldContain, "en-US")
			So(codes, ShouldContain, "en-us")
			So(codes, ShouldContain, "en")
		})

		Convey("Trim", func() {
			So(i18n.Trim("/en/abc.html"), ShouldEqual, "abc.html")
			So(i18n.Trim("/xyz.html"), ShouldEqual, "xyz.html")
		})
	})
}

func TestI18nIni(t *testing.T) {
	Convey("I18nIni", t, func() {
		i18n, err := NewI18n("en", i18nIniBytes, ".ini")
		So(err, ShouldBeNil)

		tr := i18n.Tr("meta.title")
		So(tr, ShouldEqual, "Title")

		tr = i18n.Tr("meta.xxx")
		So(tr, ShouldEqual, "meta.xxx")

		tr = i18n.Tr("a.b.c")
		So(tr, ShouldEqual, "a.b.c")

		tr = i18n.Trf("meta.link", "xyz")
		So(tr, ShouldEqual, "Linkxyz")
	})
}
