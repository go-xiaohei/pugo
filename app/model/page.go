package model

import (
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/go-xiaohei/pugo/app/parser"
	"path"
	"strings"
)

// Page contains fields for a page
type Page struct {
	Title       string  `ini:"title"`
	Slug        string  `ini:"slug"`
	Url         string  `ini:"-"`
	Permalink   string  `ini:"-"`
	HoverClass  string  `ini:"hover"`
	Template    string  `ini:"template"` // page's template for render
	Desc        string  `ini:"desc"`
	Thumb       string  `ini:"thumb"`
	Created     Time    `ini:"-"`
	Updated     Time    `ini:"-"`
	Author      *Author `ini:"-"`
	Raw         []byte
	RawType     string
	Meta        map[string]string // page's extra meta data
	ContentHTML template.HTML

	fileName string
	fileTime time.Time
}

// page's html content
func (p *Page) contentHTML() template.HTML {
	if p.RawType == "markdown" {
		return template.HTML(helper.Markdown(p.Raw))
	}
	return template.HTML(p.Raw)
}

// blocks to Page
func NewPage(blocks []parser.Block, fi os.FileInfo) (*Page, error) {
	if len(blocks) < 2 {
		return nil, ErrPostBlockError
	}
	p := &Page{
		fileName: fi.Name(),
		fileTime: fi.ModTime(),
		Meta:     make(map[string]string),
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
	if p.Template == "" {
		// default page template is page.html
		p.Template = "page.html"
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
	p.Meta = block.MapHash("meta")

	// parse markdown block
	p.RawType = blocks[1].Type()
	p.Raw = blocks[1].Bytes()

	// build url
	p.Permalink = fmt.Sprintf("/%s", p.Slug)
	p.Url = p.Permalink + ".html"

	if len(blocks) > 2 {
		// parse meta block
		block, ok := blocks[2].(parser.MetaBlock)
		if !ok {
			return nil, ErrMetaBlockWrong
		}
		p.Meta = block.MapHash("")
	}

	p.ContentHTML = p.contentHTML()
	return p, nil
}
