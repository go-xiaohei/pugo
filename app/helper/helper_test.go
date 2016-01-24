package helper

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/inconshreveable/log15.v2"
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

			So(page.PrevURL(), ShouldEqual, "ppp2")
			So(page.NextURL(), ShouldEqual, "ppp4")

			page = pager.Page(1)
			So(page.PrevURL(), ShouldEqual, "")

			page = pager.Page(20)
			So(page.NextURL(), ShouldEqual, "")
		})
	})
}

func TestMd5(t *testing.T) {
	Convey("Md5", t, func() {
		str := Md5("123456")
		So(str, ShouldEqual, "e10adc3949ba59abbe56e057f20f883e")

		str, err := Md5File("md5.go")
		So(err, ShouldBeNil)
		So(str, ShouldEqual, "d16c46931ab9d0359ad5262aa9b4a2da")
	})
}

func TestMarkdown(t *testing.T) {
	Convey("Markdown", t, func() {
		h1 := []byte("#h1")
		So(string(Markdown(h1)), ShouldEqual, `<h1 id="h1">h1</h1>`+"\n")

		code := []byte("```go\npackage main\n```")
		So(string(Markdown(code)), ShouldEqual, `<pre><code class="language-go">package main
</code></pre>
`)
	})
}

func TestLog(t *testing.T) {
	Convey("Log", t, func() {
		var buf bytes.Buffer
		l := log15.New()
		l.SetHandler(log15.StreamHandler(&buf, LogfmtFormat()))
		l.Debug("ABC|%s|%s|%s", "a", "b", "c")

		So(buf.String(), ShouldContainSubstring, "ABC|a|b|c")
	})
}
