package model

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-xiaohei/pugo/app/helper"
	"gopkg.in/ini.v1"
)

var (
	tomlPrefix         = []byte("toml")
	titleReplacer      = strings.NewReplacer(" ", "-")
	postBlockSeparator = []byte("```")
	postBriefSeparator = []byte("<!--more-->")
)

// Post contain all fields of a post content
type Post struct {
	Title      string   `toml:"title" ini:"title"`
	Slug       string   `toml:"slug" ini:"slug"`
	Desc       string   `toml:"desc" ini:"desc"`
	Date       string   `toml:"date" ini:"date"`
	Update     string   `toml:"update_date" ini:"update_date"`
	AuthorName string   `toml:"author" ini:"author"`
	Thumb      string   `toml:"thumb" ini:"thumb"`
	Draft      bool     `toml:"draft" ini:"draft"`
	TagString  []string `toml:"tags" ini:"-"`
	Tags       []*Tag   `toml:"-" ini:"-"`
	Author     *Author  `toml:"-" ini:"-"`

	dateTime   time.Time
	updateTime time.Time

	Bytes        []byte `toml:"-"`
	contentBytes []byte
	briefBytes   []byte
	postURL      string
	fileURL      string
	destURL      string
}

// SetURL set path when assemble posts
func (p *Post) SetURL(url string) {
	p.postURL = url
}

// SetDestURL set dest-url
func (p *Post) SetDestURL(url string) {
	p.destURL = url
}

// SetPlaceholder fix @placeholder in post values
func (p *Post) SetPlaceholder(stringReplacer, htmlReplacer *strings.Replacer) {
	p.Thumb = stringReplacer.Replace(p.Thumb)
	p.contentBytes = []byte(htmlReplacer.Replace(string(p.contentBytes)))
	p.briefBytes = []byte(htmlReplacer.Replace(string(p.briefBytes)))
}

// URL get url of the post
func (p *Post) URL() string {
	return p.postURL
}

// SourceURL get source file path of the post
func (p *Post) SourceURL() string {
	return filepath.ToSlash(p.fileURL)
}

// DestURL get destination file of the post after compiled
func (p *Post) DestURL() string {
	return filepath.ToSlash(p.destURL)
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

// PreviewHTML get brief html content
// deprecated
func (p *Post) PreviewHTML() template.HTML {
	return p.BriefHTML()
}

// Preview get brief html content
// deprecated
func (p *Post) Preview() []byte {
	return p.Brief()
}

// Created get create time
func (p *Post) Created() time.Time {
	return p.dateTime
}

// Updated get update time
func (p *Post) Updated() time.Time {
	return p.updateTime
}

// IsUpdated return true if updated time is not same to created time
func (p *Post) IsUpdated() bool {
	return p.updateTime.Unix() == p.dateTime.Unix()
}

func (p *Post) normalize() error {
	if p.Slug == "" {
		// use filename instead of slug, do not use title
		p.Slug = strings.TrimSuffix(filepath.Base(p.fileURL), filepath.Ext(p.fileURL))
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
	permaURL := fmt.Sprintf("/%d/%d/%d/%s", p.dateTime.Year(), p.dateTime.Month(), p.dateTime.Day(), p.Slug)
	p.postURL = permaURL + ".html"
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
	if formatType == FormatINI {
		iniObj, err := ini.Load(dataSlice[1][idx:])
		if err != nil {
			return nil, err
		}
		section := iniObj.Section("DEFAULT")
		if err = section.MapTo(post); err != nil {
			return nil, err
		}
		tagStr := section.Key("tags").Value()
		if tagStr != "" {
			post.TagString = strings.Split(tagStr, ",")
		}
		authorEmail := section.Key("author_email").Value()
		if authorEmail != "" {
			post.Author, err = newAuthorFromIniSection(section)
			if err != nil {
				return nil, err
			}
		}
	}
	post.fileURL = file
	post.Bytes = bytes.Trim(dataSlice[2], "\n")
	return post, post.normalize()
}

func parseTimeString(timeStr string) (time.Time, error) {
	timeStr = strings.TrimSpace(timeStr)
	if len(timeStr) == 0 {
		return time.Time{}, errors.New("empty time string")
	}
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
