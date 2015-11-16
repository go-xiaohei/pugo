package test

import (
	"github.com/go-xiaohei/pugo-static/model"
	"github.com/go-xiaohei/pugo-static/parser"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"testing"
)

var (
	p = parser.NewCommonParser()
)

func TestParseMeta(t *testing.T) {
	Convey("parse meta", t, func() {
		bytes, err := ioutil.ReadFile("../source/meta.md")
		So(err, ShouldBeNil)
		blocks, err := p.Parse(bytes)
		So(err, ShouldBeNil)

		Convey("check meta block", func() {
			So(blocks, ShouldHaveLength, 1)
			So(blocks[0].Type(), ShouldEqual, parser.BLOCK_INI)

			Convey("use meta block", func() {
				b, ok := blocks[0].(parser.MetaBlock)
				So(ok, ShouldBeTrue)
				So(b.Item("meta", "title"), ShouldEqual, "Pugo.Static")

				meta, err := model.NewMeta(blocks)
				So(err, ShouldBeNil)
				So(meta.Title, ShouldEqual, b.Item("meta", "title"))
			})
		})
	})
}

func TestPostMeta(t *testing.T) {
	Convey("parse post", t, func() {
		bytes, err := ioutil.ReadFile("../source/post/welcome.md")
		So(err, ShouldBeNil)
		blocks, err := p.Parse(bytes)
		So(err, ShouldBeNil)

		Convey("check post blocks", func() {
			So(blocks, ShouldHaveLength, 2)
			So(blocks[0].Type(), ShouldEqual, parser.BLOCK_INI)
			So(blocks[1].Type(), ShouldEqual, parser.BLOCK_MARKDOWN)

			Convey("use post blocks", func() {
				b, ok := blocks[0].(parser.MetaBlock)
				So(ok, ShouldBeTrue)
				So(b.Item("title"), ShouldEqual, "Welcome to Pugo.Static")

				fi, _ := os.Stat("../source/post/welcome.md")
				post, err := model.NewPost(blocks, fi)
				So(err, ShouldBeNil)
				So(post.Title, ShouldEqual, b.Item("title"))
			})
		})
	})
}

func TestPageMeta(t *testing.T) {
	Convey("parse page", t, func() {
		bytes, err := ioutil.ReadFile("../source/page/about.md")
		So(err, ShouldBeNil)
		blocks, err := p.Parse(bytes)
		So(err, ShouldBeNil)

		Convey("check page blocks", func() {
			So(blocks, ShouldHaveLength, 3)
			So(blocks[0].Type(), ShouldEqual, parser.BLOCK_INI)
			So(blocks[1].Type(), ShouldEqual, parser.BLOCK_MARKDOWN)
			So(blocks[2].Type(), ShouldEqual, parser.BLOCK_INI)

			Convey("use page blocks", func() {
				b, ok := blocks[0].(parser.MetaBlock)
				So(ok, ShouldBeTrue)
				So(b.Item("title"), ShouldEqual, "About Pugo Static")

				fi, _ := os.Stat("../source/page/about.md")
				page, err := model.NewPage(blocks, fi)
				So(err, ShouldBeNil)
				So(page.Title, ShouldEqual, b.Item("title"))
			})
		})
	})
}
