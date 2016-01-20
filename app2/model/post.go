package model

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/go-xiaohei/pugo/app2/helper"
	"github.com/naoina/toml"
)

var (
	tomlPrefix         = []byte("toml")
	titleReplacer      = strings.NewReplacer(" ", "-")
	postBlockSeparator = []byte("```")
	postBriefSeparator = []byte("<!--more-->")
	postTimeLayout     = "2006-01-02 15:04:05"
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
	Tags       []string `toml:"tags"`
	Author     *Author  `toml:"-"`

	dateTime   time.Time
	updateTime time.Time

	Bytes        []byte
	contentBytes []byte
	briefBytes   []byte
	permaURL     string
	postURL      string
}

func (p *Post) normalize() error {
	if p.Slug == "" {
		p.Slug = titleReplacer.Replace(p.Title)
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
	p.briefBytes = helper.Markdown(bytes.Split(p.Bytes, postBriefSeparator)[0])
	p.permaURL = fmt.Sprintf("/%d/%d/%d/%s", p.dateTime.Year(), p.dateTime.Month(), p.dateTime.Day(), p.Slug)
	p.postURL = p.permaURL + ".html"
	return nil
}

// NewPostOfMarkdown create new post from markdown file
func NewPostOfMarkdown(file string) (*Post, error) {
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
	post := new(Post)
	if err = toml.Unmarshal(dataSlice[1][4:], post); err != nil {
		return nil, err
	}
	post.Bytes = bytes.Trim(dataSlice[2], "\n")
	return post, post.normalize()
}
