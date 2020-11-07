package commands

import (
	"fmt"
	"pugo/pkg/builder"

	"github.com/urfave/cli/v2"
)

var (
	// Build is command of 'build'
	Build = &cli.Command{
		Name:  "build",
		Usage: "build content directory to destination",
		Flags: []cli.Flag{
			buildSourceFlag,
		},
		Action: func(ctx *cli.Context) error {
			contentDirectory := getContentFlagValue(ctx)
			// parse meta file
			metaData, err := builder.ParseMeta(contentDirectory)
			if err != nil {
				return err
			}
			fmt.Println(metaData.Meta)
			return nil
		},
	}
)
