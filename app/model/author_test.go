package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthorGroup(t *testing.T) {
	Convey("AuthorGroupEmpty", t, func() {
		ag := make(AuthorGroup, 0)
		err := ag.normalize()
		So(err, ShouldEqual, errAuthorGroupEmpty)
	})
	Convey("AuthorError", t, func() {
		author := &Author{}
		err := author.normalize()
		So(err, ShouldEqual, errAuthorInvalid)

		ag := make(AuthorGroup, 0)
		ag = append(ag, author)
		err = ag.normalize()
		So(err, ShouldEqual, errAuthorInvalid)
	})
}
