package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTree(t *testing.T) {
	tree := NewTree("")
	tree.Add("/abc.html", "abc", TreePage, 0)
	tree.Add("/abc/xyz.html", "abc-xyz", TreePage, 1)
	tree.Add("/abc/123.html", "abc-123", TreePost, 2)

	Convey("Tree", t, func() {
		children := tree.Children("abc")
		So(children, ShouldHaveLength, 2)
		So(children[0].Title, ShouldEqual, "abc-xyz")
		So(children[1].Children(), ShouldHaveLength, 0)
		So(children[1].Type, ShouldEqual, TreePost)
	})
}
