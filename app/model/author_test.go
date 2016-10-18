package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthorGroup(t *testing.T) {
	Convey("AuthorGroupEmpty", t, func() {
		ag := make(AuthorGroup, 0)
		err := ag.normalize()
		So(err.Error(), ShouldEqual, "Must add an author")
	})
}
