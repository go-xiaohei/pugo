package deploy

import "github.com/codegangsta/cli"

var (
	methods = make(map[string]Method)
)

func init() {
	Register(new(Ftp))
	Register(new(Sftp))
	Register(new(Git))
}

// Method define deploy method behavior
type Method interface {
	Create(ctx *cli.Context) (Method, error)
	Do() error
	Command() cli.Command
	String() string
}

// Register register new deploy method
func Register(m Method) {
	methods[m.String()] = m
}

// Commands get commands of all deploy methods
func Commands() []cli.Command {
	commands := make([]cli.Command, len(methods))
	i := 0
	for _, m := range methods {
		commands[i] = m.Command()
		i++
	}
	return commands
}
