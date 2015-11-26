package command

import "github.com/codegangsta/cli"

func Server() cli.Command {
	return cli.Command{
		Name:     "server",
		Usage:    "build source and server static files, watch changes for updating",
		HideHelp: true,
		Flags:    []cli.Flag{},
		Action:   serveSite(),
	}
}

func serveSite() func(ctx *cli.Context) {
	return nil
}
