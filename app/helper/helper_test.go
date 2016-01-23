package helper

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGravatar(t *testing.T) {
	Convey("Gravatar", t, func() {
		url := Gravatar("fuxiaohei@vip.qq.com", 50)
		targetURL := "https://www.gravatar.com/avatar/f72f7454ce9d710baa506394f68f4132?size=50"
		So(url, ShouldEqual, targetURL)

		url2 := Gravatar("fuxiaohei@vip.qq.com", 0)
		targetURL = "https://www.gravatar.com/avatar/f72f7454ce9d710baa506394f68f4132?size=80"
		So(url2, ShouldEqual, targetURL)
	})
}

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
)

func TestI18n(t *testing.T) {
	Convey("I18n", t, func() {
		i18n, err := NewI18n("en", i18nBytes)
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
			_, err := NewI18n("en", b)
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
	})
}
