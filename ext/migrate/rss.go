package migrate

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"

	"fmt"
	"github.com/go-xiaohei/pugo/app/builder"
	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/go-xiaohei/pugo/app/model"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/naoina/toml"
	"golang.org/x/net/html"
	"gopkg.in/inconshreveable/log15.v2"
	"net/url"
	"path"
	"time"
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
	err       error
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
func (r *RSS) Action(ctx *builder.Context) (map[string]*bytes.Buffer, error) {
	// read rss data
	resp, err := http.Get(r.Source)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		log15.Error("Migrate|RSS|Fail|%s", http.StatusText(resp.StatusCode))
		return nil, errors.New(http.StatusText(resp.StatusCode))
	}

	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)
	if err = r.parseRSSData(r.Source, buf.Bytes()); err != nil {
		return nil, err
	}
	return r.result, nil
}

func (r *RSS) parseRSSData(source string, data []byte) error {
	feed := rss.New(10, true, r.RssChannelHandler, nil)
	if err := feed.FetchBytes(source, data, nil); err != nil {
		return err
	}
	return r.err
}

func (r *RSS) RssChannelHandler(feed *rss.Feed, newChannel []*rss.Channel) {
	if len(newChannel) == 0 {
		return
	}
	cn := newChannel[0]

	log15.Debug("Migrate|RSS|Items|%d", len(cn.Items))
	// parse posts
	// it seems no way to find out whether rss item is post or page
	// so make them as post
	for _, item := range cn.Items {
		b := bytes.NewBuffer(nil)

		p := new(model.Post)
		p.Title = item.Title

		u, _ := url.Parse(item.Links[0].Href)
		slug := path.Base(u.Path)
		ext := path.Ext(slug)
		slug = strings.TrimSuffix(slug, ext)
		p.Slug = slug

		t, _ := time.Parse(time.RFC1123Z, item.PubDate)
		p.Date = t.Format("2006-01-02 15:04:05")

		if item.Author.Name != "" {
			p.AuthorName = item.Author.Name
		}

		tags := make([]string, len(item.Categories))
		for i, c := range item.Categories {
			tags[i] = c.Text
		}
		p.TagString = tags

		b2, err := toml.Marshal(*p)
		if err != nil {
			r.err = err
			return
		}
		b.WriteString("```toml\n\n")
		b.Write(b2)
		b.WriteString("\n```\n\n")

		var (
			content string
		)
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

			var contentBytes []byte
			contentBytes, r.err = helper.MarkdownBack([]byte(content))
			if r.err != nil {
				return
			}

			b.Write(contentBytes)
		}

		fileName := fmt.Sprintf("post/%s/%s.md", t.Format("2006"), slug)
		r.result[fileName] = b
		log15.Debug("RSS|Generate|%s|%d", fileName, r.result[fileName].Len())
	}

	// build meta
	// add meta.ini
	b := bytes.NewBuffer([]byte("[meta]\n"))
	fmt.Fprintf(b, `title = "%s"`+"\n", cn.Title)
	fmt.Fprintf(b, `subtitle = "%s"`+"\n", cn.SubTitle.Text)
	fmt.Fprintf(b, `keyword = ""`+"\n")
	fmt.Fprintf(b, `desc = "%s"`+"\n", cn.Description)

	var link string
	for _, ln := range cn.Links {
		if ln.Type == "" {
			link = ln.Href
			break
		}
	}
	u, _ := url.Parse(link)
	fmt.Fprintf(b, `domain = "%s"`+"\n", u.Host)
	fmt.Fprintf(b, `root = "%s"`+"\n", link)
	fmt.Fprintf(b, `cover = ""`+"\n")
	fmt.Fprintf(b, `lang = "%s"`+"\n", cn.Language)

	b.WriteString("\n")
	if cn.Author.Name == "" {
		cn.Author = rss.Author{
			Name: "PuGo",
			Uri:  "http://pugo.io",
		}
	}
	b.WriteString("[[author]]\n")
	fmt.Fprintf(b, `name = "%s"`+"\n", cn.Author.Name)
	fmt.Fprintf(b, `email = "%s"`+"\n", cn.Author.Email)
	fmt.Fprintf(b, `url = "%s"`+"\n", cn.Author.Uri)
	fmt.Fprintf(b, `bio = ""`+"\n")
	fmt.Fprintf(b, `avatar = ""`+"\n")

	b.WriteString("\n")
	// because no navigation in rss, use default value
	// comment, analytics and build setting use default text
	b.WriteString(rssDefaultNav)
	b.WriteString("\n\n")
	r.result["meta.toml"] = b
	r.result["page"] = nil
	log15.Debug("RSS|Generate|%s|%d", "meta.toml", r.result["meta.toml"].Len())
}

var (
	// default navigation for RSS
	rssDefaultNav = `
[[nav]]
link = "/"
title = "Home"
i18n = "home"
hover = "home"

[[nav]]
link = "/archive"
title = "Archive"
i18n = "archive"
hover = "archive"

[comment]
disqus = ""
duoshuo = ""

[analytics]
google = ""
baidu = ""

`
)
