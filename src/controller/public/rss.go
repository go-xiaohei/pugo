package public

import (
	"bytes"
	"github.com/go-xiaohei/pugo/src/service"
	"github.com/lunny/tango"
)

type RssController struct {
	tango.Ctx
}

func (rc *RssController) Get() {
	var buf bytes.Buffer
	if err := service.Call(service.RSS.RSS, nil, &buf); err != nil {
		panic(err)
	}
	rc.Header().Add("Content-Type", "application/rss+xml;charset=UTF-8")
	rc.Write(buf.Bytes())
}

type RobotController struct {
	tango.Ctx
}

func (rc *RobotController) Get() {
	rc.Write([]byte(`User-agent: *
Disallow: /admin/`))
}
