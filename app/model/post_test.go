package model

import (
	"fmt"
	"path"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	r  = newReplacer("")
	hr = newReplacerInHTML("")
)

func newReplacer(static string) *strings.Replacer {
	p := path.Join(static, "media")
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return strings.NewReplacer(
		"@media", p,
	)
}

func newReplacerInHTML(static string) *strings.Replacer {
	p := path.Join(static, "media")
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return strings.NewReplacer(
		`src="@media`, fmt.Sprintf(`src="%s`, p),
		`href="@media`, fmt.Sprintf(`src="%s`, p),
	)
}

func TestModelPostToml(t *testing.T) {
	Convey("ParseTomlFrontMatter", t, func() {
		p, err := NewPostOfMarkdown("testdata/post/post_toml.md", nil)
		if err != nil {
			So(err, ShouldNotBeNil)
			return
		}
		So(p.Title, ShouldEqual, "Welcome")
		So(p.Tags[0].Name, ShouldEqual, "pugo")

		t, _ := time.Parse("2006-01-02 15:04:05", "2016-03-25 12:20:20")
		So(p.Date, ShouldEqual, "2016-03-25 12:20:20")
		So(p.Created().Unix(), ShouldEqual, t.Unix())
		So(p.Updated().Format("2006-01-02"), ShouldEqual, "2016-03-26")
		So(p.IsUpdated(), ShouldEqual, true)

		So(p.Content(), ShouldHaveLength, 1768)
		So(p.Brief(), ShouldHaveLength, 1043)
		So(p.Draft, ShouldEqual, false)

		So(p.URL(), ShouldEqual, "/2016/3/25/welcome.html")
		So(p.SourceURL(), ShouldEqual, "testdata/post/post_toml.md")

		Convey("PostSetURL", func() {
			p.SetURL("/welcome.html")
			So(p.URL(), ShouldEqual, "/welcome.html")

			p.SetDestURL("dest/welcome.html")
			So(p.DestURL(), ShouldEqual, "dest/welcome.html")
		})

		Convey("PostSetPlaceholder", func() {
			p.SetPlaceholder(r, hr)
			So(p.Thumb, ShouldEqual, "/media/golang.png")
		})
	})
}

func TestModelPostIni(t *testing.T) {
	Convey("ParseIniFrontMatter", t, func() {
		p, err := NewPostOfMarkdown("testdata/post/post_ini.md", nil)
		if err != nil {
			So(err, ShouldNotBeNil)
			return
		}
		So(p.Title, ShouldEqual, "Welcome")
		So(p.Tags[0].Name, ShouldEqual, "pugo")
		So(p.Tags[1].Name, ShouldEqual, "xyz")

		So(p.Date, ShouldEqual, "2016-03-25 12:20")
		So(p.Created().Format("2006-01-02 15:04"), ShouldEqual, "2016-03-25 12:20")
		So(p.IsUpdated(), ShouldEqual, false)
		So(p.Slug, ShouldEqual, "post_ini")
	})
}

func TestModelPostMeta(t *testing.T) {
	Convey("ParsePostMeta", t, func() {
		for t, f := range ShouldPostMetaFiles() {
			file := path.Join("testdata/post", f)
			posts, err := NewPostsFrontMatter(file, t)
			So(err, ShouldBeNil)
			So(posts, ShouldContainKey, "post_toml2.md")

			p, err := NewPostOfMarkdown("testdata/post/post_toml2.md", posts["post_toml2.md"])
			So(err, ShouldBeNil)
			So(p.Content(), ShouldHaveLength, 330)
		}
	})
}

func TestModelPostWrong(t *testing.T) {
	Convey("ParsePostWrong", t, func() {
		Convey("ParseWrongTime", func() {
			_, err := NewPostOfMarkdown("testdata/post/post_wrong.md", nil)
			So(err.Error(), ShouldContainSubstring, "empty time")

			_, err = parseTimeString("")
			So(err.Error(), ShouldEqual, "empty time string")
		})

		Convey("ParseWrongFrontMatter", func() {
			_, err := NewPostOfMarkdown("testdata/post/post_wrong2.md", nil)
			So(err.Error(), ShouldContainSubstring, "unrecognized")

			_, err = NewPostOfMarkdown("testdata/post/post_wrong3.md", nil)
			So(err.Error(), ShouldContainSubstring, "need front-matter")
		})
	})
}
