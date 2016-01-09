package migrate

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	registeredTask map[string]Task

	// OutputDirectory set migrate ouput directory in global
	OutputDirectory = "source"

	ErrMigrateArgsError = fmt.Errorf("migrate args need be 'pugo migrate <type> <source>'")
	ErrMigrateUnknown   = fmt.Errorf("migrate type unknown")

	migrateMetaExtraString = `[comment]
disqus =
duoshuo =

[analytics]
google =
baidu =

[build.ignore]
-:CNAME
-:.git`
)

func init() {
	registeredTask = map[string]Task{
		TypeRSS: new(RSSTask),
	}
}

type (
	Task interface {
		Is(conf string) bool
		New(ctx *cli.Context) (Task, error)
		Type() string
		Do() (map[string]*bytes.Buffer, error)
	}
)

func Detect(ctx *cli.Context) (Task, error) {
	src := ctx.String("src")
	if len(src) == 0 || !strings.Contains(src, "://") {
		return nil, ErrMigrateArgsError
	}
	for _, task := range registeredTask {
		if task.Is(src) {
			log15.Info("Migrate.Detect.[" + task.Type() + "]")
			return task.New(ctx)
		}
	}
	return nil, ErrMigrateUnknown
}

// Register new migrate task
func Register(task Task) {
	registeredTask[task.Type()] = task
}
