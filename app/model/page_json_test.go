package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const testJSONstring = `{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "childrenAge":[19,22,28],
  "childrenHeight":[0.28,18.28,11],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "James", "last": "Murphy"},
    {"first": "Roger", "last": "Craig"}
  ],
  "billing":1289.19,
  "public":true
}`

func TestPageJson(t *testing.T) {
	j := NewJSON([]byte(testJSONstring))
	Convey("PageJson", t, func() {
		So(j.String("name.last"), ShouldEqual, "Anderson")
		So(j.String("name.last2"), ShouldEqual, "")

		So(j.Int("age"), ShouldEqual, 37)
		So(j.Int64("age"), ShouldHaveSameTypeAs, int64(37))
		So(j.Int32("age"), ShouldHaveSameTypeAs, int32(37))
		So(j.Int16("age"), ShouldHaveSameTypeAs, int16(37))
		So(j.Int8("age"), ShouldHaveSameTypeAs, int8(37))
		So(j.Int("age2"), ShouldEqual, 0)

		So(j.Float("billing"), ShouldEqual, 1289.19)
		So(j.Float32("billing"), ShouldHaveSameTypeAs, float32(1289.19))
		So(j.Float64("billing"), ShouldHaveSameTypeAs, float64(1289.19))
		So(j.Float("billing2"), ShouldEqual, 0.00)

		So(j.Bool("public"), ShouldBeTrue)
		So(j.Bool("public2"), ShouldBeFalse)

		So(j.Exist("name"), ShouldBeTrue)
		So(j.Exist("name2"), ShouldBeFalse)

		So(j.Strings("children"), ShouldHaveLength, 3)
		So(j.Strings("children"), ShouldContain, "Sara")
		So(j.Strings("children"), ShouldContain, "Alex")
		So(j.Strings("children2"), ShouldHaveLength, 0)

		So(j.Ints("childrenAge"), ShouldHaveLength, 3)
		So(j.Ints("childrenAge")[0], ShouldEqual, 19)
		So(j.Floats("childrenHeight"), ShouldHaveLength, 3)
		So(j.Floats("childrenHeight")[0], ShouldEqual, 0.28)

		So(j.Slice("children")[0].String(), ShouldEqual, "Sara")
		So(j.Get("children").Index(0).String(), ShouldEqual, "Sara")

		So(j.Map("name"), ShouldHaveLength, 2)
		So(j.Map("name")["first"].String(), ShouldEqual, "Tom")
		So(j.Get("name").Key("first").String(), ShouldEqual, "Tom")
	})
}
