package command

import (
	"github.com/go-xiaohei/pugo/ext/deploy"
	"github.com/urfave/cli"
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
	for k := range commands {
		commands[k].Flags = append(commands[k].Flags, debugFlag)
		commands[k].Before = Before
	}
	Deploy.Subcommands = commands
}
