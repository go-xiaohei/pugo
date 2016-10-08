package builder

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/go-xiaohei/pugo/app/model"
	"gopkg.in/inconshreveable/log15.v2"
)

type (
	// Source include all sources data
	Source struct {
		Meta      *model.Meta
		Nav       model.NavGroup
		Owner     *model.Author
		Authors   map[string]*model.Author
		Comment   *model.Comment
		Analytics *model.Analytics
		I18n      map[string]*helper.I18n

		Posts    model.Posts
		PostPage int
		Archive  []*model.Archive
		Pages    model.Pages
		Tags     map[string]*model.Tag
		tagPosts map[string]model.Posts
	}
)

// NewSource return new *Source from *MetaAll,
// it returns general *Source object without posts and pages,
// but meta, navigation, authors,comment and analytics are loaded.
func NewSource(all *model.MetaAll) *Source {
	s := &Source{
		Meta:      all.Meta,
		Nav:       all.NavGroup,
		Owner:     all.AuthorGroup[0],
		Comment:   all.Comment,
		Analytics: all.Analytics,
		Authors:   make(map[string]*model.Author),
	}
	for _, a := range all.AuthorGroup {
		s.Authors[a.Name] = a
	}
	return s
}

// ReadSource read source with *Context.
// parse *Context.From and read data to *Context.Source
func ReadSource(ctx *Context) {
	ctx.parseDir()
	if ctx.Err != nil {
		return
	}

	// read meta
	// then read languages,posts and pages together
	metaAll, err := ReadMeta(ctx.srcDir)
	if err != nil {
		ctx.Err = err
		return
	}
	ctx.Source = NewSource(metaAll)

	wg := helper.NewGoGroup("ReadStep")
	wg.Wrap("ReadLang", func() {
		ctx.Source.I18n = ReadLang(ctx.srcDir)
	})
	wg.Wrap("ReadPosts", func() {
		ctx.Source.Posts, ctx.Err = ReadPosts(ctx.srcDir)
	})
	wg.Wrap("ReadPages", func() {
		ctx.Source.Pages, ctx.Err = ReadPages(ctx.srcDir)
	})
	wg.Wait()
}

// ReadMeta read meta file in srcDir
func ReadMeta(srcDir string) (*model.MetaAll, error) {
	var metaFile string
	for t, f := range model.ShouldMetaFiles() {
		metaFile = filepath.Join(srcDir, f)
		if !com.IsFile(metaFile) {
			continue
		}
		log15.Debug("Read|%s", metaFile)
		bytes, err := ioutil.ReadFile(metaFile)
		if err != nil {
			return nil, err
		}
		meta, err := model.NewMetaAll(bytes, t)
		if err != nil {
			return nil, err
		}
		if meta != nil {
			return meta, nil
		}
	}
	return nil, fmt.Errorf("meta file is missing")
}

// ReadLang read languages in srcDir
func ReadLang(srcDir string) map[string]*helper.I18n {
	srcDir = filepath.Join(srcDir, "lang")
	if !com.IsDir(srcDir) {
		return nil
	}
	langs := make(map[string]*helper.I18n)
	filepath.Walk(srcDir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		ext := filepath.Ext(p)
		if ext == ".toml" || ext == ".ini" {
			log15.Debug("Read|%s", p)
			b, err := ioutil.ReadFile(p)
			if err != nil {
				log15.Warn("Read|Lang|%s|%v", p, err)
				return nil
			}
			lang := strings.TrimSuffix(filepath.Base(p), ext)
			i18n, err := helper.NewI18n(lang, b, ext)
			if err != nil {
				log15.Warn("Read|Lang|%s|%v", p, err)
				return nil
			}
			langs[lang] = i18n
		}
		return nil
	})
	return langs
}

// ReadPosts read posts files in srcDir/post
func ReadPosts(srcDir string) ([]*model.Post, error) {
	srcDir = filepath.Join(srcDir, "post")
	if !com.IsDir(srcDir) {
		return nil, fmt.Errorf("posts directory '%s' is missing", srcDir)
	}

	var posts []*model.Post
	err := filepath.Walk(srcDir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		if filepath.Ext(p) == ".md" {
			log15.Debug("Read|%s", p)
			post, err := model.NewPostOfMarkdown(p)
			if err != nil {
				log15.Warn("Read|Post|%s|%v", p, err)
			} else if post != nil && !post.Draft {
				posts = append(posts, post)
			}
			if post.Draft == true {
				log15.Warn("Draft|%s", p)
			}
		}
		return nil
	})
	sort.Sort(model.Posts(posts))
	return posts, err
}

// ReadPages read pages files in srcDir/page
func ReadPages(srcDir string) ([]*model.Page, error) {
	srcDir = filepath.Join(srcDir, "page")
	if !com.IsDir(srcDir) {
		return nil, fmt.Errorf("pages directory '%s' is missing", srcDir)
	}

	var pages []*model.Page
	err := filepath.Walk(srcDir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		if filepath.Ext(p) == ".md" {
			log15.Debug("Read|%s", p)
			rel, _ := filepath.Rel(srcDir, p)
			rel = strings.TrimSuffix(rel, filepath.Ext(rel))
			page, err := model.NewPageOfMarkdown(p, filepath.ToSlash(rel))
			if err != nil {
				log15.Warn("Read|Page|%s|%v", p, err)
			} else if page != nil && !page.Draft {
				pages = append(pages, page)
			}
			if page.Draft == true {
				log15.Warn("Draft|%s", p)
			}
		}
		return nil
	})
	return pages, err
}
