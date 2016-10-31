package model

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/go-xiaohei/pugo/app/helper"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	i18nBytes = []byte(`[meta]
title = "Title"
subtitle = "Subtitle"
link = "Link%s"

[nav]
docs = "docsabc"
guide = "guidexyz"
github = "git"`)
)

func TestMetaRead(t *testing.T) {
	Convey("ParseMeta", t, func() {
		for t, f := range ShouldMetaFiles() {
			file := path.Join("testdata", f)

			if t == FormatTOML {
				Convey("ParseMetaToml", func() {
					fileData, err := ioutil.ReadFile(file)
					if err != nil {
						So(err, ShouldBeNil)
						return
					}
					meta, err := NewMetaAll(fileData, FormatTOML)
					if err != nil {
						So(err, ShouldBeNil)
						return
					}

					So(meta.Meta.Title, ShouldEqual, "PuGo")
					So(meta.NavGroup, ShouldHaveLength, 4)
					So(meta.AuthorGroup, ShouldHaveLength, 2)
					So(meta.AuthorGroup[0].IsOwner, ShouldBeTrue)
					So(meta.Comment.IsOK(), ShouldBeTrue)

					Convey("NavGroup", func() {
						i18n, err := helper.NewI18n("en", i18nBytes, ".toml")
						So(err, ShouldBeNil)
						So(meta.NavGroup[0].Tr(i18n), ShouldEqual, "guidexyz")
						meta.NavGroup.SetPrefix("abc")
						So(meta.NavGroup[0].Link, ShouldEqual, "abc/guide.html")
						So(meta.NavGroup[0].TrLink(i18n), ShouldEqual, "/en/abc/guide.html")
						So(meta.NavGroup[3].TrLink(i18n), ShouldEqual, "https://github.com/go-xiaohei/pugo")
					})

					Convey("NavError", func() {
						ng := make(NavGroup, 0)
						ng = append(ng, &Nav{})
						err := ng.normalize()
						So(err, ShouldEqual, errNavInvalid)
					})
				})
			}

			if t == FormatINI {
				Convey("ParseMetaIni", func() {
					fileData, err := ioutil.ReadFile(file)
					if err != nil {
						So(err, ShouldBeNil)
						return
					}
					meta, err := NewMetaAll(fileData, FormatINI)
					if err != nil {
						So(err, ShouldBeNil)
						return
					}

					So(meta.Meta.Title, ShouldEqual, "PuGo")
					So(meta.NavGroup, ShouldHaveLength, 3)
					So(meta.AuthorGroup, ShouldHaveLength, 2)
					So(meta.AuthorGroup[0].IsOwner, ShouldBeTrue)
					So(meta.Comment.IsOK(), ShouldBeTrue)

					So(meta.Meta.DomainURL("/abc.html"), ShouldEqual, "http://pugo.io/docs/abc.html")
				})
			}
		}
	})

	Convey("ParseMetaError", t, func() {
		_, err := NewMetaAll([]byte("abc"), 0)
		So(err, ShouldEqual, errMetaUnsupport)

		meta := &Meta{}
		err = meta.normalize()
		So(err, ShouldEqual, errMetaInvalid)

		meta2 := &Meta{
			Title:  "pugo",
			Domain: "pugo.io",
		}
		err = meta2.normalize()
		So(err, ShouldBeNil)
	})
}
