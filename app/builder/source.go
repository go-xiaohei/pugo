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

var (
	errMetaFileMissing = fmt.Errorf("meta file is missing")
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
		Build     *model.Build
		I18n      map[string]*helper.I18n

		Posts      model.Posts
		PagePosts  map[int]*model.PagerPosts
		IndexPosts model.PagerPosts // same to PagePosts[1]
		PostPage   int
		Archive    model.Archives
		Pages      model.Pages
		Tags       map[string]*model.Tag
		TagPosts   map[string]*model.TagPosts
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
		Build:     all.Build,
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
	metaAll, err := ReadSecondMeta(ctx.srcDir)
	if err != nil {
		ctx.Err = err
		return
	}
	ctx.Source = NewSource(metaAll)

	wg := helper.NewGoGroup("ReadStep")
	wg.Wrap("ReadLang", func() error {
		ctx.Source.I18n = ReadLang(ctx.SrcLangDir())
		return nil
	})
	wg.Wrap("ReadPosts", func() error {
		if ctx.Source.Build != nil && ctx.Source.Build.DisablePost {
			return nil
		}
		var err error
		ctx.Source.Posts, err = ReadPosts(ctx)
		return err
	})
	wg.Wrap("ReadPages", func() error {
		if ctx.Source.Build != nil && ctx.Source.Build.DisablePage {
			return nil
		}
		var err error
		ctx.Source.Pages, err = ReadPages(ctx)
		return err
	})
	wg.Wait()
	if len(wg.Errors()) > 0 {
		ctx.Err = wg.Errors()[0]
	}
}

// ReadSecondMeta read meta file in srcDir
func ReadSecondMeta(srcDir string) (*model.MetaAll, error) {
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
	return nil, errMetaFileMissing
}

// ReadLang read languages in srcDir
func ReadLang(srcDir string) map[string]*helper.I18n {
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
		p = filepath.ToSlash(p)
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
func ReadPosts(ctx *Context) ([]*model.Post, error) {
	srcDir := ctx.SrcPostDir()
	if !com.IsDir(srcDir) {
		return nil, fmt.Errorf("posts directory '%s' is missing", srcDir)
	}
	log15.Info("Read|Posts|%s", srcDir)

	// try load post.toml or post.ini to read total meta file
	var (
		err      error
		postMeta = make(map[string]*model.Post)
	)
	for t, f := range model.ShouldPostMetaFiles() {
		file := filepath.Join(srcDir, f)
		if !com.IsFile(file) {
			continue
		}
		postMeta, err = model.NewPostsFrontMatter(file, t)
		if err != nil {
			return nil, err
		}
		log15.Debug("Read|PostMeta|%s", file)
		break
	}

	var posts []*model.Post
	err = filepath.Walk(srcDir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		p = filepath.ToSlash(p)
		if filepath.Ext(p) == ".md" {
			metaKey := strings.TrimPrefix(p, filepath.ToSlash(srcDir+"/"))
			log15.Debug("Read|%s|%v", p, postMeta[metaKey] != nil)
			post, err := model.NewPostOfMarkdown(p, postMeta[metaKey])
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
func ReadPages(ctx *Context) ([]*model.Page, error) {
	srcDir := ctx.SrcPageDir()
	if !com.IsDir(srcDir) {
		return nil, fmt.Errorf("pages directory '%s' is missing", srcDir)
	}
	log15.Info("Read|Pages|%s", srcDir)

	var (
		err      error
		pageMeta = make(map[string]*model.Page)
	)
	for t, f := range model.ShouldPageMetaFiles() {
		file := filepath.Join(srcDir, f)
		if !com.IsFile(file) {
			continue
		}
		pageMeta, err = model.NewPagesFrontMatter(file, t)
		if err != nil {
			return nil, err
		}
		log15.Debug("Read|PageMeta|%s", file)
		break
	}

	var pages []*model.Page
	err = filepath.Walk(srcDir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		p = filepath.ToSlash(p)
		if filepath.Ext(p) == ".md" {
			rel, _ := filepath.Rel(srcDir, p)
			rel = strings.TrimSuffix(rel, filepath.Ext(rel))
			metaKey := strings.TrimPrefix(p, filepath.ToSlash(srcDir+"/"))
			log15.Debug("Read|%s|%v", p, pageMeta[metaKey] != nil)
			page, err := model.NewPageOfMarkdown(p, filepath.ToSlash(rel), pageMeta[metaKey])
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
