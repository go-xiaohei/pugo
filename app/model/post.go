package model

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/go-xiaohei/pugo-static/app/helper"
	"github.com/go-xiaohei/pugo-static/app/parser"
	"path"
)

var (
	ErrPostBlockError = errors.New("post-block-wrong")
)

// Post contains data for a post
type Post struct {
	Title     string  `ini:"title"`
	Slug      string  `ini:"slug"`
	Permalink string  `ini:"-"`
	Url       string  `ini:"-"`
	Desc      string  `ini:"desc"` // description in a sentence
	Created   Time    `ini:"-"`
	Updated   Time    `ini:"-"`
	Author    *Author `ini:"-"`
	Tags      []*Tag  `ini:"-"`
	Raw       []byte  ``
	rawType   string

	fileName string
	fileTime time.Time
}

// post's content html
func (p *Post) ContentHTML() template.HTML {
	if p.rawType == "markdown" {
		return template.HTML(helper.Markdown(p.Raw))
	}
	return template.HTML(p.Raw)
}

// post's preview html,
// use "<!--more-->" to separate, return first part
func (p *Post) PreviewHTML() template.HTML {
	bytes := bytes.Split(p.Raw, []byte("<!--more-->"))[0]
	if p.rawType == "markdown" {
		return template.HTML(helper.Markdown(bytes))
	}
	return template.HTML(bytes)
}

// blocks  to Post
func NewPost(blocks []parser.Block, fi os.FileInfo) (*Post, error) {
	if len(blocks) != 2 {
		return nil, ErrPostBlockError
	}
	p := &Post{
		fileName: fi.Name(),
		fileTime: fi.ModTime(),
	}

	block, ok := blocks[0].(parser.MetaBlock)
	if !ok {
		return nil, ErrMetaBlockWrong
	}
	if err := block.MapTo("", p); err != nil {
		return nil, err
	}
	if p.Slug == "" {
		ext := path.Ext(fi.Name())
		p.Slug = strings.TrimSuffix(fi.Name(), ext)
	}
	tags := strings.Split(block.Item("tags"), ",")
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if t != "" {
			p.Tags = append(p.Tags, NewTag(t))
		}
	}

	p.Created = NewTime(block.Item("date"), p.fileTime)
	p.Updated = p.Created
	if upStr := block.Item("update_date"); upStr != "" {
		p.Updated = NewTime(upStr, p.fileTime)
	}
	p.Author = &Author{
		Name:  block.Item("author"),
		Email: block.Item("author_email"),
		Url:   block.Item("author_url"),
	}

	// parse markdown block
	p.rawType = blocks[1].Type()
	p.Raw = blocks[1].Bytes()

	// build url
	p.Permalink = fmt.Sprintf("/%d/%d/%d/%s", p.Created.Year, p.Created.Month, p.Created.Day, p.Slug)
	p.Url = p.Permalink + ".html"

	return p, nil
}

// posts list
type Posts []*Post

// implement sort.Sort interface
func (p Posts) Len() int           { return len(p) }
func (p Posts) Less(i, j int) bool { return p[i].Created.Raw.Unix() > p[j].Created.Raw.Unix() }
func (p Posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// archive stat for posts
type Archive struct {
	Year  int // each list by year
	Posts []*Post
}

// posts to archive
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
