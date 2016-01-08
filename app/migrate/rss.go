package migrate

import (
	"bytes"
	"fmt"
	"html"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/helper"
	rss "github.com/jteeuwen/go-pkg-rss"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	TypeRSS = "RSS"
)

var (
	_ Task = new(RSSTask)

	ErrRSSSchemaWrong = fmt.Errorf("Migrate RSS need schema 'rss+http://' or 'rss+https://'")

	purlRSSEncodeBeginTag = "<http://purl.org/rss/1.0/modules/content/:encoded>"
	purlRSSEncodeEndTag   = "</http://purl.org/rss/1.0/modules/content/:encoded>"
)

type (
	RSSTask struct {
		opt    *RSSOption
		result map[string]*bytes.Buffer
	}
	RSSOption struct {
		Dest     string
		Source   string
		IsRemote bool
	}
)

func (rs *RSSTask) Is(conf string) bool {
	return strings.HasPrefix(conf, "rss+")
}

func (rs *RSSTask) New(ctx *cli.Context) (Task, error) {
	u, err := url.Parse(ctx.String("src"))
	if err != nil {
		return nil, err
	}

	opt := &RSSOption{
		Dest: ctx.String("dest"),
	}
	if len(u.Scheme) <= 4 {
		return nil, ErrRSSSchemaWrong
	}
	// get real schema, to get remote rss source
	u.Scheme = u.Scheme[4:]
	opt.IsRemote = true
	opt.Source = u.String()

	return &RSSTask{
		opt: opt,
	}, nil
}

func (rs *RSSTask) Type() string {
	return TypeRSS
}

func (rs *RSSTask) Dir() string {
	return "dir"
}

func (rs *RSSTask) Do() (map[string]*bytes.Buffer, error) {
	rs.result = make(map[string]*bytes.Buffer)

	feed := rss.New(10, true, rs.chanHandler, nil)
	log15.Debug("RSS.Read." + rs.opt.Source)
	if err := feed.Fetch(rs.opt.Source, nil); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %s: %s\n", rs.opt.Source, err)
		return nil, err
	}
	return rs.result, nil
}

func (rs *RSSTask) chanHandler(feed *rss.Feed, newChannel []*rss.Channel) {
	if len(newChannel) == 0 {
		return
	}
	cn := newChannel[0]

	// parse posts
	// it seems no way to find out whether rss item is post or page
	// so make them as post
	for _, item := range cn.Items {
		b := bytes.NewBuffer(nil)
		b.WriteString("```ini\n\n")

		b.WriteString(fmt.Sprintf("title = %s\n\n", item.Title))

		u, _ := url.Parse(item.Links[0].Href)
		slug := path.Base(u.Path)
		ext := path.Ext(slug)
		slug = strings.TrimSuffix(slug, ext)
		b.WriteString(fmt.Sprintf("slug = %s\n\n", slug))
		b.WriteString(fmt.Sprintf("desc = \n\n"))

		t, _ := time.Parse(time.RFC1123Z, item.PubDate)
		b.WriteString(fmt.Sprintf("date = %s\n\n", t.Format("2006-01-02 15:04:05")))

		if item.Author.Name != "" {
			b.WriteString(fmt.Sprintf("author = %s\n\n", item.Author.Name))
		}

		tags := make([]string, len(item.Categories))
		for i, c := range item.Categories {
			tags[i] = c.Text
		}
		b.WriteString(fmt.Sprintf("tags = %s\n\n", strings.Join(tags, ",")))

		b.WriteString("```\n\n")

		var content string
		if item.Content != nil {
			content = item.Content.Text
		}
		if content == "" {
			content = item.Description
		}
		if content != "" {
			if strings.HasPrefix(content, purlRSSEncodeBeginTag) {
				content = strings.TrimPrefix(content, purlRSSEncodeBeginTag)
			}
			if strings.HasSuffix(content, purlRSSEncodeEndTag) {
				content = strings.TrimSuffix(content, purlRSSEncodeEndTag)
			}
			content = html.UnescapeString(content)
			content = helper.HTML2Markdown(content)
			b.WriteString(content)
		}

		fileName := fmt.Sprintf("post/%s/%s.md", t.Format("2006"), slug)
		rs.result[fileName] = b
	}

	// fmt.Println(rs.result)
}
