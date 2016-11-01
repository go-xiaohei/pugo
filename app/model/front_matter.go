package model

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/ini.v1"
)

// NewPostsFrontMatter parse post meta file to create post data
func NewPostsFrontMatter(file string, t FormatType) (map[string]*Post, error) {
	metas := make(map[string]*Post)

	if t == FormatTOML {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		if err = toml.Unmarshal(data, &metas); err != nil {
			return nil, err
		}
	}

	if t == FormatINI {
		iniObj, err := ini.Load(file)
		if err != nil {
			return nil, err
		}
		for _, s := range iniObj.SectionStrings() {
			s2 := strings.Trim(s, `"`)
			if s2 == "DEFAULT" {
				continue
			}
			post := new(Post)
			if err = newPostFromIniSection(iniObj.Section(filepath.ToSlash(s)), post); err != nil {
				return nil, err
			}
			metas[s2] = post
		}
	}
	return metas, nil
}

// NewPagesFrontMatter parse page meta file to create page data
func NewPagesFrontMatter(file string, t FormatType) (map[string]*Page, error) {
	metas := make(map[string]*Page)

	if t == FormatTOML {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		if err = toml.Unmarshal(data, &metas); err != nil {
			return nil, err
		}
	}

	if t == FormatINI {
		iniObj, err := ini.Load(file)
		if err != nil {
			return nil, err
		}
		for _, s := range iniObj.SectionStrings() {
			s2 := strings.Trim(s, `"`)
			if s2 == "DEFAULT" {
				continue
			}
			if strings.HasSuffix(s2, ".meta") {
				continue
			}
			page := new(Page)
			if err = newPageFromIniObject(iniObj, page, s, `"`+s2+`.meta"`); err != nil {
				return nil, err
			}
			metas[s2] = page
		}
	}

	for k, page := range metas {
		if page.Node {
			page.Slug = k
			if err := page.normalize(); err != nil {
				return nil, err
			}
		}
	}

	return metas, nil
}

func newPostFromIniSection(section *ini.Section, post *Post) error {
	var err error
	if err = section.MapTo(post); err != nil {
		return err
	}
	tagStr := section.Key("tags").Value()
	if tagStr != "" {
		post.TagString = strings.Split(tagStr, ",")
	}
	authorEmail := section.Key("author_email").Value()
	if authorEmail != "" {
		post.Author, err = newAuthorFromIniSection(section)
		if err != nil {
			return err
		}
	}
	return nil
}

func newPageFromIniObject(iniObj *ini.File, page *Page, sectionName, metaSectionName string) error {
	var err error
	section := iniObj.Section(sectionName)
	if err = section.MapTo(page); err != nil {
		return err
	}
	authorEmail := section.Key("author_email").Value()
	if authorEmail != "" {
		page.Author, err = newAuthorFromIniSection(section)
		if err != nil {
			return err
		}
	}
	metaData := iniObj.Section(metaSectionName).KeysHash()
	if len(metaData) > 0 {
		page.Meta = make(map[string]interface{})
		for k, v := range metaData {
			page.Meta[k] = v
		}
	}
	return nil
}
