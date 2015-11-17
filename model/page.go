package model

import (
	"fmt"
	"github.com/go-xiaohei/pugo-static/parser"
	"html/template"
	"os"
	"time"
)

type Page struct {
	Title      string `ini:"title"`
	Slug       string `ini:"slug"`
	Url        string `ini:"-"`
	HoverClass string `ini:"hover"`
	Template   string `ini:"template"`
	Desc       string `ini:"desc"`
	Created    Time   `ini:"-"`
	Updated    Time   `ini:"-"`
	Author     Author `ini:"-"`
	Raw        []byte
	rawType    string
	Meta       map[string]string

	fileName string
	fileTime time.Time
}

func (p *Page) ContentHTML() template.HTML {
	if p.rawType == "markdown" {
		return template.HTML(markdown(p.Raw))
	}
	return template.HTML(p.Raw)
}

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
	if p.Template == "" {
		p.Template = "page.html"
	}

	p.Created = NewTime(block.Item("date"), p.fileTime)
	p.Updated = p.Created
	if upStr := block.Item("update_date"); upStr != "" {
		p.Updated = NewTime(upStr, p.fileTime)
	}
	p.Author = Author{
		Name:  block.Item("author"),
		Email: block.Item("author_email"),
		Url:   block.Item("author_url"),
	}
	p.Meta = block.MapHash("meta")

	// parse markdown block
	p.rawType = blocks[1].Type()
	p.Raw = blocks[1].Bytes()

	// build url
	p.Url = fmt.Sprintf("/%s", p.Slug)

	if len(blocks) > 2 {
		// parse meta block
		block, ok := blocks[2].(parser.MetaBlock)
		if !ok {
			return nil, ErrMetaBlockWrong
		}
		p.Meta = block.MapHash("")
	}
	return p, nil
}
