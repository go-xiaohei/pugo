package public

import (
	"bytes"
	"github.com/fuxiaohei/pugo/src/service"
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
	rc.Req().Header.Add("Content-Type", "application/rss+xml;charset=UTF-8")
	rc.Write(buf.Bytes())
}
