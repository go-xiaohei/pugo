package model

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/go-xiaohei/pugo/app/parser"
)

// Page contains fields for a page
type Page struct {
	Title     string `ini:"title"`
	Slug      string `ini:"slug"`
	URL       string `ini:"-"`
	Permalink string `ini:"-"`

	HoverClass string `ini:"hover"`
	Template   string `ini:"template"` // page's template for render
	Desc       string `ini:"desc"`
	Thumb      string `ini:"thumb"`
	Lang       string `ini:"lang"` // language

	Created Time    `ini:"-"`
	Updated Time    `ini:"-"`
	Author  *Author `ini:"-"`

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
		return helper.Bytes2MarkdownHTML(p.Raw)
	}
	return template.HTML(p.Raw)
}

// NewPage parses blocks to Page
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
		URL:   block.Item("author_url"),
	}
	if p.Author.Name == "" {
		p.Author = nil
	}
	p.Meta = block.MapHash("meta")

	// parse markdown block
	p.RawType = blocks[1].Type()
	p.Raw = blocks[1].Bytes()

	// build url
	p.Permalink = fmt.Sprintf("/%s", p.Slug)
	p.URL = p.Permalink + ".html"
	if p.Lang != "" {
		// use language url
		p.Permalink = fmt.Sprintf("/%s%s", strings.ToLower(p.Lang), p.Permalink)
		p.URL = fmt.Sprintf("/%s%s", strings.ToLower(p.Lang), p.URL)
	} else {
		p.Lang = "-"
	}

	p.ContentHTML = p.contentHTML()
	return p, nil
}

type (
	// PageNodeGroup defins page nodes
	PageNodeGroup map[string]map[string]*pageNode
	pageNode      struct {
		URL       string
		Permalink string
	}
)

// NewPageNodeGroup generates page nodes from pages
func NewPageNodeGroup(pages []*Page) PageNodeGroup {
	m := make(map[string]map[string]*pageNode)
	for _, page := range pages {
		if len(m[page.Slug]) == 0 {
			m[page.Slug] = make(map[string]*pageNode)
		}
		m[page.Slug][page.Lang] = &pageNode{
			URL:       page.URL,
			Permalink: page.Permalink,
		}
	}
	return PageNodeGroup(m)
}

// URL returns node url by slug and language
func (png PageNodeGroup) URL(slug string, lang string) string {
	languages := helper.NewI18nLanguageCode(lang)
	for _, l := range languages {
		if url, ok := png[slug][l]; ok {
			return url.URL
		}
	}
	if url := png[slug]["-"]; url != nil {
		return url.URL
	}
	return ""
}

// Permalink returns node permalink by slug and language
func (png PageNodeGroup) Permalink(slug string, lang string) string {
	languages := helper.NewI18nLanguageCode(lang)
	for _, l := range languages {
		if url, ok := png[slug][l]; ok {
			return url.Permalink
		}
	}
	if url := png[slug]["-"]; url != nil {
		return url.Permalink
	}
	return ""
}
