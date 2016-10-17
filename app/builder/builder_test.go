package builder

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-xiaohei/pugo/app/helper"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/urfave/cli"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/inconshreveable/log15.v2/ext"
)

func init() {
	log15.Root().SetHandler(log15.LvlFilterHandler(log15.LvlDebug, ext.FatalHandler(log15.StreamHandler(os.Stderr, helper.LogfmtFormat()))))
}

func TestBuildSimple(t *testing.T) {
	log15.Root().SetHandler(log15.StreamHandler(ioutil.Discard, helper.LogfmtFormat()))
	Convey("Build Simple", t, func() {
		ctx := NewContext(&cli.Context{}, "../../source", "../../dest", "../../theme/default")
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
	Convey("Valid Context", t, func() {
		ctx := NewContext(&cli.Context{}, "../../source2", "../../dest", "../../theme/default")
		ShouldBeFalse(ctx.IsValid())
		ShouldNotBeNil(ctx.Err)

		ctx = NewContext(&cli.Context{}, "", "", "")
		ShouldBeFalse(ctx.IsValid())
	})

	Convey("Context Directory", t, func() {
		ctx := NewContext(&cli.Context{}, "../../source", "../../dest", "../../theme/default")
		ShouldEqual(ctx.SrcDir(), "../source")
		ShouldEqual(ctx.DstDir(), "../dest")
	})

	Convey("Context Scheme", t, func() {
		ctx := NewContext(&cli.Context{}, "dir://../../source", "dir://../../dest", "../../theme/default")
		ShouldBeTrue(ctx.IsValid())

		ctx = NewContext(&cli.Context{}, "ftp://../../source", "ftp://../../dest", "../../theme/default")
		ShouldBeTrue(ctx.IsValid())
		ShouldNotBeNil(ctx.Err)
	})
}
