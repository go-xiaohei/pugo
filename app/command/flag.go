package command

import "github.com/codegangsta/cli"

var (
	buildToFlag = cli.StringFlag{
		Name:  "to",
		Value: "dir://public",
		Usage: "write files to destination or remote path",
	}
	buildFromFlag = cli.StringFlag{
		Name:  "from",
		Value: "dir://source",
		Usage: "read files from source directory or remote path",
	}
	themeFlag = cli.StringFlag{
		Name:  "theme",
		Value: "dir://theme/default",
		Usage: "theme to use (located in flag directory)",
	}
	debugFlag = cli.BoolFlag{
		Name:  "debug",
		Usage: "print more logs in debug mode",
	}
	watchFlag = cli.BoolFlag{
		Name:  "watch",
		Usage: "watch changes and rebuild files",
	}
	addrFlag = cli.StringFlag{
		Name:  "addr",
		Value: "0.0.0.0:9899",
		Usage: "http server address",
	}
	newToFlag = cli.StringFlag{
		Name:  "to",
		Value: "dir://source",
		Usage: "create new content to this directory",
	}
)
