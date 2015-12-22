package deploy

import (
	"errors"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/builder"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	registeredDeployWay      map[string]DeployTask
	ErrDeployConfFormatError = errors.New("deploy format need be type:conf_string")
	ErrDeployUnknown         = errors.New("deploy way is unknown")
)

func init() {
	registeredDeployWay = map[string]DeployTask{
		TYPE_GIT:  new(GitTask),
		TYPE_FTP:  new(FtpTask),
		TYPE_SFTP: new(SftpTask),
	}
}

type (
	// DeployTask defines the methods of a deploy task
	DeployTask interface {
		Is(conf string) bool                               // is this deploy task
		New(conf string) (DeployTask, error)               // new instance
		Type() string                                      // task type name
		Dir() string                                       // the build target directory for the deployment
		Do(b *builder.Builder, ctx *builder.Context) error // deploy logic
	}
)

// Detect the deploy task to run
func Detect(ctx *cli.Context) (DeployTask, error) {
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
