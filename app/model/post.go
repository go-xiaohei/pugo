package model

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/helper"
	"golang.org/x/net/html"
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
	Title      string       `toml:"title" ini:"title"`
	Slug       string       `toml:"slug" ini:"slug"`
	Desc       string       `toml:"desc" ini:"desc"`
	Date       string       `toml:"date" ini:"date"`
	Update     string       `toml:"update_date" ini:"update_date"`
	AuthorName string       `toml:"author" ini:"author"`
	Thumb      string       `toml:"thumb" ini:"thumb"`
	Draft      bool         `toml:"draft" ini:"draft"`
	TagString  []string     `toml:"tags" ini:"-"`
	Tags       []*Tag       `toml:"-" ini:"-"`
	Author     *Author      `toml:"-" ini:"-"`
	Index      []*PostIndex `toml:"-" ini:"-"`

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
	return p.updateTime.Unix() != p.dateTime.Unix()
}

func (p *Post) normalize() error {
	if p.Slug == "" {
		// use filename instead of slug, do not use title
		p.Slug = strings.TrimSuffix(filepath.Base(p.fileURL), filepath.Ext(p.fileURL))
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
	p.briefBytes = helper.Markdown(bytes.Split(p.Bytes, postBriefSeparator)[0])
	permaURL := fmt.Sprintf("/%d/%d/%d/%s", p.dateTime.Year(), p.dateTime.Month(), p.dateTime.Day(), p.Slug)
	p.postURL = permaURL + ".html"
	for _, t := range p.TagString {
		p.Tags = append(p.Tags, NewTag(t))
	}
	p.Index = newPostIndexs(bytes.NewReader(p.contentBytes))
	return nil
}

// NewPostOfMarkdown create new post from markdown file
func NewPostOfMarkdown(file string, post *Post) (*Post, error) {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if len(fileBytes) < 3 {
		return nil, fmt.Errorf("post content is too less")
	}

	if post == nil {
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

		post = new(Post)
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
			if err = newPostFromIniSection(section, post); err != nil {
				return nil, err
			}
		}
		post.Bytes = bytes.Trim(dataSlice[2], "\n")
	} else {
		post.Bytes = bytes.Trim(fileBytes, "\n")
	}
	post.fileURL = file
	if post.Date == "" {
		t, _ := com.FileMTime(file)
		post.dateTime = time.Unix(t, 0)
	}
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

// PostIndex is index of post
type PostIndex struct {
	Level    int
	Title    string
	Archor   string
	Children []*PostIndex
	Link     string
	Parent   *PostIndex
}

// Print prints post indexs friendly
func (p *PostIndex) Print() {
	fmt.Println(strings.Repeat("#", p.Level), p)
	for _, c := range p.Children {
		c.Print()
	}
}

func newPostIndexs(r io.Reader) []*PostIndex {
	var (
		z = html.NewTokenizer(r)

		currentLevel    int
		currentText     string
		currentLinkText string
		currentArchor   string
		nodeDeep        int

		indexs []*PostIndex
	)
	for {
		token := z.Next()
		if token == html.ErrorToken {
			break
		}
		if token == html.EndTagToken {
			if nodeDeep == 1 && currentLevel > 0 {
				indexs = append(indexs, &PostIndex{
					Level:  currentLevel,
					Title:  currentText,
					Link:   currentLinkText,
					Archor: currentArchor,
				})
				currentLevel = 0
				currentText = ""
				currentLinkText = ""
				currentArchor = ""
			}
			nodeDeep--
			continue
		}
		if token == html.StartTagToken {
			name, hasAttr := z.TagName()
			lv := parsePostIndexLevel(name)

			if lv > 0 {
				currentLevel = lv
				if hasAttr {
					for {
						k, v, isMore := z.TagAttr()
						if bytes.Equal(k, []byte("id")) {
							currentArchor = string(v)
						}
						if !isMore {
							break
						}
					}
				}
			}
			nodeDeep++

			if currentLevel > 0 && string(name) == "a" {
				if hasAttr {
					for {
						k, v, isMore := z.TagAttr()
						if bytes.Equal(k, []byte("href")) {
							currentLinkText = string(v)
						}
						if !isMore {
							break
						}
					}
				}
			}
		}
		if token == html.TextToken && currentLevel > 0 {
			currentText += string(z.Text())
		}
	}
	indexs = assemblePostIndex(indexs)
	return indexs
}

func assemblePostIndex(indexList []*PostIndex) []*PostIndex {
	var (
		list    []*PostIndex
		lastIdx int
		lastN   *PostIndex
	)
	for i, n := range indexList {
		if i == 0 {
			list = append(list, n)
			lastIdx = 0
			continue
		}
		lastN = list[lastIdx]
		if lastN.Level < n.Level {
			n.Parent = lastN
			lastN.Children = append(lastN.Children, n)
		} else {
			list = append(list, n)
			lastIdx++
		}
	}
	for _, n := range list {
		if len(n.Children) > 1 {
			n.Children = assemblePostIndex(n.Children)
		}
	}
	return list
}

func parsePostIndexLevel(name []byte) int {
	if bytes.Equal(name, []byte("h1")) {
		return 1
	}
	if bytes.Equal(name, []byte("h2")) {
		return 2
	}
	if bytes.Equal(name, []byte("h3")) {
		return 3
	}
	if bytes.Equal(name, []byte("h4")) {
		return 4
	}
	if bytes.Equal(name, []byte("h5")) {
		return 5
	}
	if bytes.Equal(name, []byte("h6")) {
		return 6
	}
	return 0
}
