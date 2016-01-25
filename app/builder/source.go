package builder

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
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

		Posts    []*model.Post
		PostPage int
		Archive  []*model.Archive
		Pages    []*model.Page
		Tags     map[string]*model.Tag
		tagPosts map[string][]*model.Post
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
	var (
		srcDir  = ""
		destDir = ""
	)
	if srcDir, ctx.Err = toDir(ctx.From); ctx.Err != nil {
		return
	}
	if !com.IsDir(srcDir) {
		ctx.Err = fmt.Errorf("Directory '%s' is missing", srcDir)
		return
	}
	ctx.srcDir = srcDir
	log15.Info("Build|Source|%s", srcDir)

	if destDir, ctx.Err = toDir(ctx.To); ctx.Err != nil {
		return
	}
	ctx.dstDir = destDir

	// read meta
	// then read posts,
	// then read pages
	metaAll, err := ReadMeta(srcDir)
	if err != nil {
		ctx.Err = err
		return
	}
	ctx.Source = NewSource(metaAll)

	if ctx.Source.I18n, ctx.Err = ReadLang(srcDir); ctx.Err != nil {
		return
	}
	if ctx.Source.Posts, ctx.Err = ReadPosts(srcDir); ctx.Err != nil {
		return
	}
	if ctx.Source.Pages, ctx.Err = ReadPages(srcDir); ctx.Err != nil {
		return
	}
}

// ReadMeta read meta.toml in srcDir
func ReadMeta(srcDir string) (*model.MetaAll, error) {
	metaFile := filepath.Join(srcDir, "meta.toml")
	if !com.IsFile(metaFile) {
		return nil, fmt.Errorf("Meta.toml is missing")
	}
	bytes, err := ioutil.ReadFile(metaFile)
	if err != nil {
		return nil, err
	}
	meta := &model.MetaAll{}
	if err = toml.Unmarshal(bytes, meta); err != nil {
		return nil, err
	}
	if err = meta.Normalize(); err != nil {
		return nil, err
	}
	log15.Debug("Build|Load|%s", metaFile)
	return meta, nil
}

// ReadLang read languages in srcDir
func ReadLang(srcDir string) (map[string]*helper.I18n, error) {
	srcDir = filepath.Join(srcDir, "lang")
	if !com.IsDir(srcDir) {
		return nil, nil
	}
	langs := make(map[string]*helper.I18n)
	err := filepath.Walk(srcDir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		if filepath.Ext(p) == ".toml" {
			log15.Debug("Build|Load|%s", p)
			b, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}
			lang := strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
			i18n, err := helper.NewI18n(lang, b)
			if err != nil {
				return err
			}
			langs[lang] = i18n
		}
		return nil
	})
	return langs, err
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
			log15.Debug("Build|Load|%s", p)
			post, err := model.NewPostOfMarkdown(p)
			if err != nil {
				return err
			}
			posts = append(posts, post)
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
			log15.Debug("Build|Load|%s", p)
			page, err := model.NewPageOfMarkdown(p)
			if err != nil {
				return err
			}
			pages = append(pages, page)
		}
		return nil
	})
	return pages, err
}

func toDir(urlString string) (string, error) {
	if !strings.Contains(urlString, "://") {
		return urlString, nil
	}
	if strings.HasPrefix(urlString, "dir://") {
		return strings.TrimPrefix(urlString, "dir://"), nil
	}
	return "", errors.New("Directory need schema dir://")
}