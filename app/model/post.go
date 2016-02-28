package model

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-xiaohei/pugo/app/helper"
)

var (
	tomlPrefix         = []byte("toml")
	titleReplacer      = strings.NewReplacer(" ", "-")
	postBlockSeparator = []byte("```")
	postBriefSeparator = []byte("<!--more-->")
)

// Post contain all fields of a post content
type Post struct {
	Title      string   `toml:"title"`
	Slug       string   `toml:"slug"`
	Desc       string   `toml:"desc"`
	Date       string   `toml:"date"`
	Update     string   `toml:"update_date"`
	AuthorName string   `toml:"author"`
	Thumb      string   `toml:"thumb"`
	TagString  []string `toml:"tags"`
	Tags       []*Tag   `toml:"-"`
	Author     *Author  `toml:"-"`

	dateTime   time.Time
	updateTime time.Time

	Bytes        []byte `toml:"-"`
	contentBytes []byte
	briefBytes   []byte
	permaURL     string
	postURL      string
	treeURL      string
}

// FixURL fix path when assemble posts
func (p *Post) FixURL(prefix string) {
	p.permaURL = path.Join(prefix, p.permaURL)
	p.postURL = path.Join(prefix, p.postURL)
}

// FixPlaceholder fix @placeholder in post values
func (p *Post) FixPlaceholder(r, hr *strings.Replacer) {
	p.Thumb = r.Replace(p.Thumb)
	p.contentBytes = []byte(hr.Replace(string(p.contentBytes)))
	p.briefBytes = []byte(hr.Replace(string(p.briefBytes)))
}

// TreeURL get tree path of the post, use to create *Tree
func (p *Post) TreeURL() string {
	return p.treeURL
}

// URL get url of the post
func (p *Post) URL() string {
	return p.postURL
}

// Permalink get permalink of the post
func (p *Post) Permalink() string {
	return p.permaURL
}

// ContentHTML get html content
func (p *Post) ContentHTML() template.HTML {
	return template.HTML(p.contentBytes)
}

// Content get html content bytes
func (p *Post) Content() []byte {
	return p.contentBytes
}

// BriefHTML get brief html content
func (p *Post) BriefHTML() template.HTML {
	return template.HTML(p.briefBytes)
}

// Brief get brief content bytes
func (p *Post) Brief() []byte {
	return p.briefBytes
}

// Created get create time
func (p *Post) Created() time.Time {
	return p.dateTime
}

// Updated get update time
func (p *Post) Updated() time.Time {
	return p.updateTime
}

func (p *Post) normalize() error {
	if p.Slug == "" {
		p.Slug = titleReplacer.Replace(p.Title)
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
	p.briefBytes = helper.Markdown(bytes.Split(p.Bytes, postBriefSeparator)[0])
	p.permaURL = fmt.Sprintf("/%d/%d/%d/%s", p.dateTime.Year(), p.dateTime.Month(), p.dateTime.Day(), p.Slug)
	p.postURL = p.permaURL + ".html"
	p.treeURL = p.permaURL
	for _, t := range p.TagString {
		p.Tags = append(p.Tags, NewTag(t))
	}
	return nil
}

// NewPostOfMarkdown create new post from markdown file
func NewPostOfMarkdown(file string) (*Post, error) {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if len(fileBytes) < 3 {
		return nil, fmt.Errorf("post content is too less")
	}
	dataSlice := bytes.SplitN(fileBytes, postBlockSeparator, 3)
	if len(dataSlice) != 3 {
		return nil, fmt.Errorf("post need front-matter block and markdown block")
	}

	idx := getFirstBreakByte(dataSlice[1])
	if idx == 0 {
		return nil, fmt.Errorf("post need front-matter block and markdown block")
	}

	formatType := detectFormat(string(dataSlice[1][:idx]))
	if formatType == 0 {
		return nil, fmt.Errorf("post front-matter block is unrecognized")
	}

	post := new(Post)
	if formatType == FormatTOML {
		if err = toml.Unmarshal(dataSlice[1][idx:], post); err != nil {
			return nil, err
		}
	}
	post.Bytes = bytes.Trim(dataSlice[2], "\n")
	return post, post.normalize()
}

// Posts are posts list
type Posts []*Post

// implement sort.Sort interface
func (p Posts) Len() int           { return len(p) }
func (p Posts) Less(i, j int) bool { return p[i].dateTime.Unix() > p[j].dateTime.Unix() }
func (p Posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func parseTimeString(timeStr string) (time.Time, error) {
	timeStr = strings.TrimSpace(timeStr)
	if len(timeStr) == len("2006-01-02") {
		return time.Parse("2006-01-02", timeStr)
	}
	if len(timeStr) == len("2006-01-02 15:04") {
		return time.Parse("2006-01-02 15:04", timeStr)
	}
	if len(timeStr) == len("2006-01-02 15:04:05") {
		return time.Parse("2006-01-02 15:04:05", timeStr)
	}
	return time.Time{}, errors.New("unknown time string")
}

func getFirstBreakByte(data []byte) int {
	for i, v := range data {
		if v == 10 {
			return i
		}
	}
	return 0
}
