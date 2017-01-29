package helper

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
