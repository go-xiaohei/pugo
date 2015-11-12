package model

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-xiaohei/pugo-static/parser"
	"gopkg.in/ini.v1"
	"html/template"
	"os"
	"strings"
	"time"
)

var (
	ErrPostBlockError = errors.New("post-block-wrong")
)

type Post struct {
	Title   string `ini:"title"`
	Slug    string `ini:"slug"`
	Url     string `ini:"-"`
	Desc    string `ini:"desc"` // description in a sentence
	Created Time   `ini:"-"`
	Updated Time   `ini:"-"`
	Author  Author `ini:"-"`
	Tags    []Tag  `ini:"-"`
	Raw     []byte ``
	rawType string

	fileName string
	fileTime time.Time
}

func (p *Post) ContentHTML() template.HTML {
	if p.rawType == "markdown" {
		return template.HTML(markdown(p.Raw))
	}
	return template.HTML(p.Raw)
}

func (p *Post) PreviewHTML() template.HTML {
	bytes := bytes.Split(p.Raw, []byte("<!--more-->"))[0]
	if p.rawType == "markdown" {
		return template.HTML(markdown(bytes))
	}
	return template.HTML(bytes)
}

func NewPost(blocks []parser.Block, fi os.FileInfo) (*Post, error) {
	if len(blocks) != 2 {
		return nil, ErrPostBlockError
	}
	p := &Post{
		fileName: fi.Name(),
		fileTime: fi.ModTime(),
	}

	// parse first ini block
	iniF, err := ini.Load(blocks[0].Bytes())
	if err != nil {
		return nil, err
	}
	section := iniF.Section("DEFAULT")
	if err := section.MapTo(p); err != nil {
		return nil, err
	}
	tags := section.Key("tags").Strings(",")
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if t != "" {
			p.Tags = append(p.Tags, NewTag(t))
		}
	}

	p.Created = NewTime(section.Key("date").String(), p.fileTime)
	p.Updated = p.Created
	if upStr := section.Key("update_date").String(); upStr != "" {
		p.Updated = NewTime(upStr, p.fileTime)
	}
	p.Author = Author{
		Name:  section.Key("author").String(),
		Email: section.Key("author_email").String(),
		Url:   section.Key("author_url").String(),
	}

	// parse markdown block
	p.rawType = blocks[1].Type()
	p.Raw = blocks[1].Bytes()

	// build url
	p.Url = fmt.Sprintf("/%d/%d/%d/%s", p.Created.Year, p.Created.Month, p.Created.Day, p.Slug)
	return p, nil
}

type Posts []*Post

func (p Posts) Len() int {
	return len(p)
}

func (p Posts) Less(i, j int) bool {
	return p[i].Created.Raw.Unix() > p[j].Created.Raw.Unix()
}

func (p Posts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type Archive struct {
	Year  int
	Posts []*Post
}

func NewArchive(posts []*Post) []*Archive {
	archives := []*Archive{}
	var (
		last, lastYear int
	)
	for _, p := range posts {
		if len(archives) == 0 {
			archives = append(archives, &Archive{
				Year:  p.Created.Year,
				Posts: []*Post{p},
			})
			continue
		}
		last = len(archives) - 1
		lastYear = archives[last].Year
		if lastYear == p.Created.Year {
			archives[last].Posts = append(archives[last].Posts, p)
			continue
		}
		archives = append(archives, &Archive{
			Year:  p.Created.Year,
			Posts: []*Post{p},
		})
	}
	return archives
}
