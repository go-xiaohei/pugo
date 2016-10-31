package model

import (
	"fmt"
	"path"
	"testing"
	"time"

	"github.com/Unknwon/com"
	. "github.com/smartystreets/goconvey/convey"
)

func TestModelPageToml(t *testing.T) {
	Convey("ParseTomlFrontMatter", t, func() {
		p, err := NewPageOfMarkdown("testdata/page/page_toml.md", "page/page_toml", nil)
		if err != nil {
			So(err, ShouldNotBeNil)
			return
		}
		So(p.Title, ShouldEqual, "Welcome")

		t, _ := time.Parse("2006-01-02 15:04:05", "2016-03-25 12:20:20")
		So(p.Date, ShouldEqual, "2016-03-25 12:20:20")
		So(p.Created().Unix(), ShouldEqual, t.Unix())
		So(p.Updated().Format("2006-01-02"), ShouldEqual, "2016-03-26")
		So(p.IsUpdated(), ShouldEqual, true)

		So(p.Content(), ShouldHaveLength, 1768)
		So(p.Draft, ShouldEqual, false)

		So(p.URL(), ShouldEqual, "/page/page_toml.html")
		So(p.SourceURL(), ShouldEqual, "testdata/page/page_toml.md")

		So(p.Meta, ShouldContainKey, "key")

		Convey("PageSetURL", func() {
			p.SetURL("/welcome.html")
			So(p.URL(), ShouldEqual, "/welcome.html")

			p.SetDestURL("dest/welcome.html")
			So(p.DestURL(), ShouldEqual, "dest/welcome.html")
		})

	})
}

func TestModelPageIni(t *testing.T) {
	Convey("ParseIniFrontMatter", t, func() {
		p, err := NewPageOfMarkdown("testdata/page/page_ini.md", "page/page_ini", nil)
		if err != nil {
			So(err, ShouldNotBeNil)
			return
		}
		So(p.Title, ShouldEqual, "Welcome")

		t, _ := com.FileMTime("testdata/page/page_ini.md")
		So(p.Created().Unix(), ShouldEqual, t)
		So(p.IsUpdated(), ShouldEqual, false)
		So(p.Meta, ShouldContainKey, "key")
	})
}

func TestModelPageMeta(t *testing.T) {
	Convey("ParsePageMeta", t, func() {
		for t, f := range ShouldPageMetaFiles() {
			file := path.Join("testdata/page", f)
			pages, err := NewPagesFrontMatter(file, t)
			if t == FormatINI {
				So(err, ShouldBeNil)
				So(pages, ShouldContainKey, "page_ini.md")

				p, err := NewPageOfMarkdown("testdata/page/page_ini.md", "page/page_ini", pages["page_ini.md"])
				So(err, ShouldBeNil)
				So(p.Content(), ShouldHaveLength, 2136)
				So(p.Meta, ShouldContainKey, "key2")
			}
			if t == FormatTOML {
				So(err, ShouldBeNil)
				So(pages, ShouldContainKey, "page_toml.md")

				p, err := NewPageOfMarkdown("testdata/page/page_toml.md", "page/page_toml", pages["page_toml.md"])
				So(err, ShouldBeNil)
				So(p.Content(), ShouldHaveLength, 2110)
				So(p.Meta, ShouldContainKey, "key3")

				Convey("PageNodeInMeta", func() {
					So(pages, ShouldContainKey, "page/node")
					So(pages["page/node"].Node, ShouldBeTrue)
					So(pages["page/node"].Content(), ShouldHaveLength, 0)
				})
			}
		}
	})
}

func TestModePageNode(t *testing.T) {
	Convey("ParsePageNode", t, func() {
		p, err := NewPageOfMarkdown("testdata/page/page_node.md", "page/page_node", nil)
		if err != nil {
			So(err, ShouldNotBeNil)
			return
		}
		So(p.Title, ShouldEqual, "PageNode")
		fmt.Println(string(p.Content()))
		So(p.Content(), ShouldHaveLength, 0)
	})
}

func TestModelPageWrong(t *testing.T) {
	Convey("ParseWrongFrontMatter", t, func() {
		// page can parse post data
		_, err := NewPageOfMarkdown("testdata/post/post_wrong2.md", "post_wrong2.md", nil)
		So(err.Error(), ShouldContainSubstring, "unrecognized")

		_, err = NewPageOfMarkdown("testdata/post/post_wrong3.md", "post_wrong3.md", nil)
		So(err.Error(), ShouldContainSubstring, "need front-matter")

		_, err = NewPageOfMarkdown("testdata/page/page_wrong.md", "page_wrong.md", nil)
		So(err.Error(), ShouldContainSubstring, "page content is too less")
	})
}
