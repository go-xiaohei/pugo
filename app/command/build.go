package command

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/builder"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	// Build is command of 'build'
	Build = cli.Command{
		Name:  "build",
		Usage: "build static files",
		Flags: []cli.Flag{
			buildFromFlag,
			buildToFlag,
			themeFlag,
		},
		Action: build,
	}
)

func build(c *cli.Context) {
	ctx := builder.NewContext(
		c.String("from"),
		c.String("to"),
		c.String("theme"),
	)
	if !ctx.IsValid() {
		log15.Crit("Build|must have values in 'from', 'to' & 'theme'")
	}
	fmt.Println(ctx)
	builder.Build(ctx)
}
