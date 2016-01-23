package model

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/naoina/toml"
	"html/template"
)

// Page contain all fields of a page content
type Page struct {
	Title      string                 `toml:"title"`
	Slug       string                 `toml:"slug"`
	Desc       string                 `toml:"desc"`
	Date       string                 `toml:"date"`
	Update     string                 `toml:"update_date"`
	AuthorName string                 `toml:"author"`
	NavHover   string                 `toml:"hover"`
	Template   string                 `toml:"template"`
	Lang       string                 `toml:"lang"`
	Bytes      []byte                 `toml:"-"`
	Meta       map[string]interface{} `toml:"meta"`
	Sort       int                    `toml:"sort"`
	Author     *Author                `toml:"-"`

	permaURL     string
	pageURL      string
	treeURL      string
	contentBytes []byte
	dateTime     time.Time
	updateTime   time.Time
}

func (p *Page) TreeURL() string {
	return p.treeURL
}

func (p *Page) URL() string {
	return p.pageURL
}

func (p *Page) Permalink() string {
	return p.permaURL
}

func (p *Page) ContentHTML() template.HTML {
	return template.HTML(p.contentBytes)
}

func (p *Page) Content() []byte {
	return p.contentBytes
}

// FixURL fix path when assemble posts
func (p *Page) FixURL(prefix string) {
	p.permaURL = path.Join(prefix, p.permaURL)
	p.pageURL = path.Join(prefix, p.pageURL)
}

// FixPlaceholder fix @placeholder in post values
func (p *Page) FixPlaceholder(hr *strings.Replacer) {
	p.contentBytes = []byte(hr.Replace(string(p.contentBytes)))
}

// Created get create time
func (p *Page) Created() time.Time {
	return p.dateTime
}

// Updated get update time
func (p *Page) Updated() time.Time {
	return p.updateTime
}

func (p *Page) normalize() error {
	if p.Slug == "" {
		p.Slug = titleReplacer.Replace(p.Title)
	}
	if p.Template == "" {
		p.Template = "page.html"
	}
	var err error
	if p.dateTime, err = time.Parse(postTimeLayout, p.Date); err != nil {
		return err
	}
	if p.Update == "" {
		p.Update = p.Date
		p.updateTime = p.dateTime
	} else {
		if p.updateTime, err = time.Parse(postTimeLayout, p.Update); err != nil {
			return err
		}
	}
	p.contentBytes = helper.Markdown(p.Bytes)
	if p.Lang == "" {
		p.permaURL = fmt.Sprintf("/%s", p.Slug)
		p.pageURL = p.permaURL + ".html"
	} else {
		p.permaURL = fmt.Sprintf("/%s/%s", p.Lang, p.Slug)
		p.pageURL = p.permaURL + ".html"
	}
	p.treeURL = p.Slug
	return nil
}

// NewPageOfMarkdown create new page from markdown file
func NewPageOfMarkdown(file string) (*Page, error) {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	dataSlice := bytes.SplitN(fileBytes, postBlockSeparator, 3)
	if len(dataSlice) != 3 {
		return nil, fmt.Errorf("post need toml block and markdown block")
	}
	if !bytes.HasPrefix(dataSlice[1], tomlPrefix) {
		return nil, fmt.Errorf("post need toml block at first")
	}
	page := new(Page)
	if err = toml.Unmarshal(dataSlice[1][4:], page); err != nil {
		return nil, err
	}
	page.Bytes = bytes.Trim(dataSlice[2], "\n")
	return page, page.normalize()
}
