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
	Title   string
	Slug    string
	Url     string
	Desc    string // description in a sentence
	Created Time
	Updated Time
	Author  Author
	Tags    []Tag
	Raw     []byte
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
	p.Title = section.Key("title").String()
	p.Slug = section.Key("slug").String()
	p.Desc = section.Key("desc").String()
	tags := section.Key("tags").Strings(",")
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if t != "" {
			p.Tags = append(p.Tags, NewTag(t))
		}
	}

	ct, err := time.Parse("2006-01-02", section.Key("date").String())
	if err != nil {
		return nil, err
	}
	p.Created = NewTime(ct)
	if upStr := section.Key("update_date").String(); upStr != "" {
		ut, err := time.Parse("2006-01-02", upStr)
		if err != nil {
			return nil, err
		}
		p.Updated = NewTime(ut)
	} else {
		p.Updated = p.Created
		// p.Updated = NewTime(p.fileTime)
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
