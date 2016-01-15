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

	"errors"
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/helper"
	rss "github.com/jteeuwen/go-pkg-rss"
	"gopkg.in/inconshreveable/log15.v2"
	"io"
	"net/http"
)

const (
	// TypeRSS is type string of RSS migration
	TypeRSS = "RSS"
)

var (
	_ Task = new(RSSTask)

	// ErrRSSSchemaWrong is error of wrong schema to run rss migration
	ErrRSSSchemaWrong = fmt.Errorf("Migrate RSS need schema 'rss+http://' or 'rss+https://'")

	purlRSSEncodeBeginTag = "<http://purl.org/rss/1.0/modules/content/:encoded>"
	purlRSSEncodeEndTag   = "</http://purl.org/rss/1.0/modules/content/:encoded>"
)

type (
	// RSSTask is a migration of RSS
	RSSTask struct {
		opt    *RSSOption
		result map[string]*bytes.Buffer
		err    error
	}
	// RSSOption is option of RSSTask
	RSSOption struct {
		Dest     string
		Source   string
		IsRemote bool
	}
)

// Is checks conf string is in RSSTask
func (rs *RSSTask) Is(conf string) bool {
	return strings.HasPrefix(conf, "rss+")
}

// New returns new RSSTask with cli.Context,
// src need be RSSTask schema,
// dest set output directory
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

// Type gets RSSTask's type string
func (rs *RSSTask) Type() string {
	return TypeRSS
}

// Do does migration of RSS,
func (rs *RSSTask) Do() (map[string]*bytes.Buffer, error) {
	rs.result = make(map[string]*bytes.Buffer)

	feed := rss.New(10, true, rs.chanHandler, nil)
	log15.Debug("RSS.Read." + rs.opt.Source)

	r, err := http.Get(rs.opt.Source)
	if err != nil {
		return nil, err
	}
	if r.StatusCode >= 400 {
		log15.Error("RSS.Read.Fail", "status", r.StatusCode)
		return nil, errors.New(http.StatusText(r.StatusCode))
	}

	var buf bytes.Buffer
	io.Copy(&buf, r.Body)
	if err := feed.FetchBytes(rs.opt.Source, buf.Bytes(), nil); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %s: %s\n", rs.opt.Source, err)
		return nil, err
	}
	return rs.result, rs.err
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
			// todo : try to download media file
			content, rs.err = helper.HTML2Markdown(content)
			if rs.err != nil {
				return
			}
			b.WriteString(content)
		}

		fileName := fmt.Sprintf("post/%s/%s.md", t.Format("2006"), slug)
		rs.result[fileName] = b
		log15.Debug("RSS.Generate.[" + fileName + "]")
	}

	// add meta.ini
	b := bytes.NewBuffer([]byte("\n"))
	b.WriteString(fmt.Sprintf("title = %s\n\n", cn.Title))
	b.WriteString(fmt.Sprintf("subtitle = %s\n\n", cn.SubTitle.Text))
	b.WriteString(fmt.Sprintf("keyword = \n\n"))
	b.WriteString(fmt.Sprintf("desc = %s\n\n", cn.Description))

	var link string
	for _, ln := range cn.Links {
		if ln.Type == "" {
			link = ln.Href
			break
		}
	}
	u, _ := url.Parse(link)
	b.WriteString(fmt.Sprintf("domain = %s\n\n", u.Host))
	b.WriteString(fmt.Sprintf("root = %s\n\n", link))
	b.WriteString("cover = \n\n")
	b.WriteString(fmt.Sprintf("lang = %s\n\n", cn.Language))

	// because no navigation in rss, use default value
	b.WriteString(rssDefaultNav)
	b.WriteString("\n\n")

	if cn.Author.Name != "" {
		b.WriteString("[author]\n-:" + cn.Author.Name + "\n\n")
		b.WriteString("[author." + cn.Author.Name + "]\n")
		b.WriteString(fmt.Sprintf("name = %s\n", cn.Author.Name))
		b.WriteString(fmt.Sprintf("email = %s\n", cn.Author.Email))
		b.WriteString(fmt.Sprintf("url = %s\n", cn.Author.Uri))
		b.WriteString("avatar = \nbio = \n")
	}

	// comment, analytics and build setting use default text
	b.WriteString("\n")
	b.WriteString(migrateMetaExtraString)
	b.WriteString("\n")

	rs.result["meta.ini"] = b
	rs.result["page/about.md"] = nil // use nil value to create empty directory
	log15.Debug("RSS.Generate.[meta.ini]")
}

var (
	// default navigation for RSS
	rssDefaultNav = `[nav]
-:home
-:archive

[nav.home]
link = /
title = Home
i18n = home
hover = home

[nav.archive]
link = /archive
title = Archive
i18n = archive
hover = archive`
)
