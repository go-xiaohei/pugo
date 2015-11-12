package model

import (
	"fmt"
	"github.com/go-xiaohei/pugo-static/parser"
	"gopkg.in/ini.v1"
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

	// parse first ini block
	iniF, err := ini.Load(blocks[0].Bytes())
	if err != nil {
		return nil, err
	}
	section := iniF.Section("DEFAULT")
	if err := section.MapTo(p); err != nil {
		return nil, err
	}
	if p.Template == "" {
		p.Template = "page.html"
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
	p.Url = fmt.Sprintf("/%s", p.Slug)

	if len(blocks) > 2 {
		// parse meta block
		iniF, err = ini.Load(blocks[2].Bytes())
		if err != nil {
			return nil, err
		}
		p.Meta = iniF.Section("DEFAULT").KeysHash()
	}
	return p, nil
}
