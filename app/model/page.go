package model

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/helper"
	"gopkg.in/ini.v1"
)

// Page contain all fields of a page content
type Page struct {
	Title      string                 `toml:"title" ini:"title"`
	Slug       string                 `toml:"slug" ini:"slug"`
	Desc       string                 `toml:"desc" ini:"desc"`
	Date       string                 `toml:"date" ini:"date"`
	Update     string                 `toml:"update_date" ini:"update_date"`
	AuthorName string                 `toml:"author" ini:"author"`
	NavHover   string                 `toml:"hover" ini:"hover"`
	Template   string                 `toml:"template" ini:"template"`
	Lang       string                 `toml:"lang" ini:"lang"`
	Bytes      []byte                 `toml:"-"`
	Meta       map[string]interface{} `toml:"meta" ini:"-"`
	Sort       int                    `toml:"sort" ini:"sort"`
	Author     *Author                `toml:"-" ini:"-"`
	Draft      bool                   `toml:"draft" ini:"draft"`
	Node       bool                   `toml:"node" ini:"node"`

	pageURL      string
	fileURL      string
	destURL      string
	contentBytes []byte
	dateTime     time.Time
	updateTime   time.Time
}

// DestURL is dest url of node
func (p *Page) DestURL() string {
	return p.destURL
}

// URL is page's url
func (p *Page) URL() string {
	return p.pageURL
}

// SourceURL get source file path of the page
func (p *Page) SourceURL() string {
	return filepath.ToSlash(p.fileURL)
}

// ContentHTML is page's content html
func (p *Page) ContentHTML() template.HTML {
	return template.HTML(p.contentBytes)
}

// Content is page's content bytes
func (p *Page) Content() []byte {
	return p.contentBytes
}

// SetURL set path when assemble posts
func (p *Page) SetURL(url string) {
	p.pageURL = url
}

// SetDestURL set path when assemble posts
func (p *Page) SetDestURL(url string) {
	p.destURL = url
}

// SetPlaceholder fix @placeholder in post values
func (p *Page) SetPlaceholder(htmlReplacer *strings.Replacer) {
	p.contentBytes = []byte(htmlReplacer.Replace(string(p.contentBytes)))
}

// Created get create time
func (p *Page) Created() time.Time {
	return p.dateTime
}

// Updated get update time
func (p *Page) Updated() time.Time {
	return p.updateTime
}

// IsUpdated return true if updated time is not same to created time
func (p *Page) IsUpdated() bool {
	return p.updateTime.Unix() != p.dateTime.Unix()
}

func (p *Page) normalize() error {
	if p.Template == "" {
		p.Template = "page.html"
	}
	var err error
	if p.Date != "" {
		if p.dateTime, err = parseTimeString(p.Date); err != nil {
			return err
		}
	}
	if p.Update == "" {
		p.Update = p.Date
		p.updateTime = p.dateTime
	} else {
		if p.updateTime, err = parseTimeString(p.Update); err != nil {
			return err
		}
	}
	p.contentBytes = helper.Markdown(p.Bytes)
	p.pageURL = "/" + p.Slug
	if !p.Node {
		p.pageURL = fmt.Sprintf("/%s", p.Slug) + ".html"
	}
	return nil
}

// NewPageOfMarkdown create new page from markdown file
func NewPageOfMarkdown(file, slug string, page *Page) (*Page, error) {
	// page-node need not read file
	if page != nil && page.Node == true {
		return page, nil
	}
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if len(fileBytes) < 3 {
		return nil, fmt.Errorf("page content is too less")
	}
	if page == nil {
		dataSlice := bytes.SplitN(fileBytes, postBlockSeparator, 3)
		if len(dataSlice) != 3 {
			return nil, fmt.Errorf("page need front-matter block and markdown block")
		}

		idx := getFirstBreakByte(dataSlice[1])
		if idx == 0 {
			return nil, fmt.Errorf("page need front-matter block and markdown block")
		}

		formatType := detectFormat(string(dataSlice[1][:idx]))
		if formatType == 0 {
			return nil, fmt.Errorf("page front-matter block is unrecognized")
		}

		page = new(Page)
		if formatType == FormatTOML {
			if err = toml.Unmarshal(dataSlice[1][idx:], page); err != nil {
				return nil, err
			}
		}
		if formatType == FormatINI {
			iniObj, err := ini.Load(dataSlice[1][idx:])
			if err != nil {
				return nil, err
			}
			if err = newPageFromIniObject(iniObj, page, "DEFAULT", "meta"); err != nil {
				return nil, err
			}
		}
		if page.Node == false {
			page.Bytes = bytes.Trim(dataSlice[2], "\n")
		}
	} else {
		page.Bytes = bytes.Trim(fileBytes, "\n")
	}
	page.fileURL = file
	if page.Slug == "" {
		page.Slug = slug
	}
	if page.Date == "" && page.Node == false { // page-node need not time
		t, _ := com.FileMTime(file)
		page.dateTime = time.Unix(t, 0)
	}
	return page, page.normalize()
}

// Pages means pages list
type Pages []*Page

// BySlug return a page by slug string
func (p Pages) BySlug(slug string) *Page {
	for _, page := range p {
		if page.Slug == slug {
			return page
		}
	}
	return nil
}
