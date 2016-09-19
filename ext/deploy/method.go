package deploy

import "github.com/urfave/cli"

var (
	methods = make(map[string]Method)
)

func init() {
	Register(new(Ftp), new(Sftp), new(Git))
	Register(new(Qiniu), new(AwsS3))
}

// Method define deploy method behavior
type Method interface {
	Create(ctx *cli.Context) (Method, error)
	Do() error
	Command() cli.Command
	String() string
}

// Register register new deploy method
func Register(ms ...Method) {
	for _, m := range ms {
		methods[m.String()] = m
	}
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
