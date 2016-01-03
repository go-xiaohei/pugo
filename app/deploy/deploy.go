package deploy

import (
	"errors"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/builder"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	registeredDeployWay map[string]Task

	// ErrDeployConfFormatError means deploy task conf string is wrong
	ErrDeployConfFormatError = errors.New("deploy format need be type:conf_string")
	// ErrDeployUnknown means unknown deploy task way
	ErrDeployUnknown = errors.New("deploy way is unknown")
)

func init() {
	registeredDeployWay = map[string]Task{
		TypeGit:  new(GitTask),
		TypeFtp:  new(FtpTask),
		TypeSftp: new(SftpTask),
	}
}

type (
	// Task defines the methods of a deploy task
	Task interface {
		Is(conf string) bool                               // is this deploy task
		New(conf string) (Task, error)                     // new instance
		Type() string                                      // task type name
		Dir() string                                       // the build target directory for the deployment
		Do(b *builder.Builder, ctx *builder.Context) error // deploy logic
	}
)

// Detect the deploy task to run
func Detect(ctx *cli.Context) (Task, error) {
	// need protocol separator
	conf := ctx.String("dest")
	if !strings.Contains(conf, "://") {
		return nil, nil
	}
	// use all ways' objects to detect
	for _, way := range registeredDeployWay {
		if way.Is(conf) {
			log15.Info("Deploy.Detect.[" + way.Type() + "]")
			return way.New(conf)
		}
	}
	return nil, ErrDeployUnknown
}

// Register new Deploy Task
func Register(task Task) {
	registeredDeployWay[task.Type()] = task
}
