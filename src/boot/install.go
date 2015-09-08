package boot
import "github.com/codegangsta/cli"

var (
    installCommand = cli.Command{
        Name:   "install",
        Usage:  "install pugo blog service with default configurations and data",
        Action: Install,
    }
)

func Install(ctx *cli.Context) {

}