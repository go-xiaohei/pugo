package migrate

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-xiaohei/pugo/app/builder"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	rssScheme = []string{
		"rss+http://",
		"rss+https://",
	}
	purlRSSEncodeBeginTag = "<http://purl.org/rss/1.0/modules/content/:encoded>"
	purlRSSEncodeEndTag   = "</http://purl.org/rss/1.0/modules/content/:encoded>"
)

// RSS migrate contents from rss source
type RSS struct {
	Directory string
	Source    string
	result    map[string]*bytes.Buffer
}

// Name return "RSS"
func (r *RSS) Name() string {
	return "RSS"
}

// Detect detect proper Task
func (r *RSS) Detect(ctx *builder.Context) (Task, error) {
	for _, prefix := range rssScheme {
		if strings.HasPrefix(ctx.From, prefix) {
			source := strings.TrimPrefix(ctx.From, "rss+")
			log15.Debug("Migrate|RSS|%s", source)
			ctx.From = "dir://source"
			log15.Debug("Migrate|RSS|To|%s", ctx.From)
			return &RSS{
				Directory: ctx.SrcDir(),
				Source:    source,
				result:    make(map[string]*bytes.Buffer),
			}, nil
		}
	}
	return nil, nil
}

// Action do rss migration to source
func (r *RSS) Action(ctx *builder.Context) error {
	// read rss data
	resp, err := http.Get(r.Source)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		log15.Error("Migrate|RSS|Fail|%s", http.StatusText(resp.StatusCode))
		return errors.New(http.StatusText(resp.StatusCode))
	}

	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)
	return r.parseRSSData(buf.Bytes())
}

func (r *RSS) parseRSSData(data []byte) error {
	return nil
}
