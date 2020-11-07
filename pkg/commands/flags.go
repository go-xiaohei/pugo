package commands

import (
	"github.com/urfave/cli/v2"
)

var (
	buildSourceFlag = &cli.StringFlag{
		Name:  "source",
		Value: "source",
		Usage: "raw content directory that build to website, same as 'content' flag",
	}
	buildContentFlag = &cli.StringFlag{
		Name:  "content",
		Value: "source",
		Usage: "raw content directory that build to website",
	}
)

func getContentFlagValue(ctx *cli.Context) string {
	content := ctx.String("content")
	if content == "" {
		return ctx.String("source")
	}
	return content
}
