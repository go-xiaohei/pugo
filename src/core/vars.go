package core
import "github.com/codegangsta/cli"

const (
    PUGO_VERSION = "1.0"
    PUGO_VERSION_STATE = "alpha"
    PUGO_VERSION_DATE = "20150910"
    PUGO_NAME = "Pugo"
    PUGO_DESCRIPTION = "a simple golang blog engine"
    PUGO_AUTHOR = "fuxiaohei"
    PUGO_AUTHOR_EMAIL = "fuxiaohei@vip.qq.com"

    RUM_MODE = "prod" // prod || debug
)

var(
    App *cli.App
)