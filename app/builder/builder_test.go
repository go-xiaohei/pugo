package builder

import (
	"io/ioutil"
	"testing"

	"github.com/go-xiaohei/pugo/app/helper"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/urfave/cli"
	"gopkg.in/inconshreveable/log15.v2"
)

func init() {
	//log15.Root().SetHandler(log15.LvlFilterHandler(log15.LvlDebug, ext.FatalHandler(log15.StreamHandler(os.Stderr, helper.LogfmtFormat()))))
	log15.Root().SetHandler(log15.StreamHandler(ioutil.Discard, helper.LogfmtFormat()))
}

func TestBuildSimple(t *testing.T) {
	Convey("Build Simple", t, func() {
		ctx := NewContext(&cli.Context{}, "../../source", "../../dest", "../../source/theme/default")
		ShouldBeTrue(ctx.IsValid())
		ShouldBeNil(ctx.Err)

		Convey("Build All", func() {
			Build(ctx)
			ShouldBeNil(ctx.Err)
		})

		Convey("Read All", func() {
			Read(ctx)
			ShouldBeNil(ctx.Err)
		})
	})
}

func TestBuildContext(t *testing.T) {
	Convey("Invalid Context", t, func() {
		ctx := NewContext(&cli.Context{}, "../../source2", "../../dest", "../../source/theme/default")
		ShouldBeFalse(ctx.IsValid())
		ShouldNotBeNil(ctx.Err)

		ctx = NewContext(&cli.Context{}, "", "", "")
		ShouldBeFalse(ctx.IsValid())
	})

	Convey("Context Directory", t, func() {
		ctx := NewContext(&cli.Context{}, "../../source", "../../dest", "../../source/theme/default")
		ShouldEqual(ctx.SrcDir(), "../source")
		ShouldEqual(ctx.DstDir(), "../dest")
		ShouldEqual(ctx.SrcPostDir(), "../source/post")
		ShouldEqual(ctx.SrcPageDir(), "../source/page")
		ShouldEqual(ctx.SrcMediaDir(), "../source/media")
		ShouldEqual(ctx.SrcLangDir(), "../source/lang")
	})
}
