package model

import (
	"io/ioutil"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
}
