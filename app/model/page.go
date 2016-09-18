package model

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
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

	pageURL      string
	treeURL      string
	fileURL      string
	contentBytes []byte
	dateTime     time.Time
	updateTime   time.Time
}

// TreeURL is tree url of node
func (p *Page) TreeURL() string {
	return p.treeURL
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

// FixURL fix path when assemble posts
func (p *Page) FixURL(prefix string) {
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
	if p.Template == "" {
		p.Template = "page.html"
	}
	var err error
	if p.dateTime, err = parseTimeString(p.Date); err != nil {
		return err
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
	p.pageURL = fmt.Sprintf("/%s", p.Slug) + ".html"
	p.treeURL = p.Slug
	return nil
}

// NewPageOfMarkdown create new page from markdown file
func NewPageOfMarkdown(file, slug string) (*Page, error) {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if len(fileBytes) < 3 {
		return nil, fmt.Errorf("page content is too less")
	}
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

	page := new(Page)
	page.fileURL = file
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
		section := iniObj.Section("DEFAULT")
		if err = section.MapTo(page); err != nil {
			return nil, err
		}
		authorEmail := section.Key("author_email").Value()
		if authorEmail != "" {
			page.Author, err = newAuthorFromIniSection(section)
			if err != nil {
				return nil, err
			}
		}
		metaData := iniObj.Section("meta").KeysHash()
		if len(metaData) > 0 {
			page.Meta = make(map[string]interface{})
			for k, v := range metaData {
				page.Meta[k] = v
			}
		}
	}
	if page.Slug == "" {
		page.Slug = slug
	}
	page.Bytes = bytes.Trim(dataSlice[2], "\n")
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
