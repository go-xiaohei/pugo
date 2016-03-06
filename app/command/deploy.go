package command

import (
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/ext/deploy"
)

var (
	// Deploy is command `deploy`
	Deploy = cli.Command{
		Name:  "deploy",
		Usage: "deploy static website",
	}
)

func init() {
	commands := deploy.Commands()
	for k, c := range commands {
		c.Flags = append(c.Flags, debugFlag)
		c.Before = Before
		commands[k] = c
	}
	Deploy.Subcommands = commands
}
