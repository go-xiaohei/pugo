package command

import "github.com/urfave/cli"

var (
	buildDestFlag = cli.StringFlag{
		Name:  "dest",
		Value: "dest",
		Usage: "write files to destination directory",
	}
	buildSourceFlag = cli.StringFlag{
		Name:  "source",
		Value: "source",
		Usage: "read files from source directory",
	}
	buildThemeFlag = cli.StringFlag{
		Name:  "theme",
		Value: "theme/default",
		Usage: "theme to use (located in flag directory)",
	}
	debugFlag = cli.BoolFlag{
		Name:  "debug",
		Usage: "print more logs in debug mode",
	}
	buildWatchFlag = cli.BoolFlag{
		Name:  "watch",
		Usage: "watch changes and rebuild files",
	}
	addrFlag = cli.StringFlag{
		Name:  "addr",
		Value: "0.0.0.0:9899",
		Usage: "http server address",
	}
	serveStaticFlag = cli.BoolFlag{
		Name:  "static",
		Usage: "just serve static file, no build",
	}
	newToFlag = cli.StringFlag{
		Name:  "to",
		Value: "dir://source",
		Usage: "create new content to this directory",
	}
	newOnlyDocFlag = cli.BoolFlag{
		Name:  "doc",
		Usage: "extract documentation data",
	}
)
