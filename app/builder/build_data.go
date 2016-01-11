package builder

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/go-xiaohei/pugo/app/model"
	"github.com/go-xiaohei/pugo/app/parser"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	//ErrParserMissing means that it can't detect proper parser
	ErrParserMissing = errors.New("Parser-Unknown")
)

// ReadData parses data to context, if parsed all data to context for renders
func (b *Builder) ReadData(ctx *Context) {
	if b.readMeta(ctx); ctx.Error != nil {
		return
	}

	// load theme after meta data reading finished
	b.render.SetFunc("url", func(str ...string) string {
		if len(str) > 0 {
			if ur, _ := url.Parse(str[0]); ur != nil {
				if ur.Host != "" {
					return str[0]
				}
			}
		}
		return path.Join(append([]string{ctx.Meta.Base}, str...)...)
	})
	b.render.SetFunc("fullUrl", func(str ...string) string {
		return ctx.Meta.Root + path.Join(str...)
	})

	// post meta process
	if b.afterMeta(ctx); ctx.Error != nil {
		return
	}

	// load contents
	if b.readContents(ctx); ctx.Error != nil {
		return
	}
}

// read meta data, from meta.md,nav.md and comment.md
// they are global values.
func (b *Builder) readMeta(ctx *Context) {
	bytes, err := ioutil.ReadFile(path.Join(b.opt.SrcDir, "meta.ini"))
	if err != nil {
		ctx.Error = err
		return
	}
	fileBytes := []byte("```ini\n")
	fileBytes = append(fileBytes, bytes...)
	fileBytes = append(fileBytes, []byte("\n```\n")...)
	blocks, err := b.parseBytes(fileBytes)
	if err != nil {
		ctx.Error = err
		return
	}

	total, err := model.NewAllMeta(blocks)
	if err != nil {
		ctx.Error = err
		return
	}
	ctx.Meta = total.Meta
	ctx.Navs = total.Nav
	ctx.Authors = total.Authors
	ctx.Comment = total.Comment
	ctx.Conf = total.Conf
	ctx.Analytics = total.Analytics
}

// do works after meta data,
// generate proper path, link and owner
func (b *Builder) afterMeta(ctx *Context) {
	// read theme
	// if error, do not need to load contents
	theme, err := b.render.Load(b.opt.Theme)
	if err != nil {
		ctx.Error = err
		return
	}
	ctx.Theme = theme

	// assign copy directory
	staticDir := ctx.Theme.Static()
	ctx.mediaPath = path.Join(ctx.Meta.Base, path.Base(staticDir), path.Base(b.opt.MediaDir))
	ctx.staticPath = path.Join(ctx.Meta.Base, path.Base(staticDir))

	replacer := replaceGlobalVars(b, ctx)

	ctx.Meta.Cover = string(replacer([]byte(ctx.Meta.Cover)))

	// read i18n data
	i18nGroup, err := model.NewI18nGroup(path.Join(b.opt.SrcDir, "lang"))
	if err != nil {
		ctx.Error = err
		return
	}
	ctx.I18nGroup = i18nGroup
	if ctx.Meta.Lang != "" {
		ctx.I18n = ctx.I18nGroup.Find(ctx.Meta.Lang)
	}
	if ctx.I18n == nil {
		log15.Warn("Lang." + ctx.Meta.Lang + ".Missing")
		ctx.I18n = helper.NewI18nEmpty()
	} else {
		log15.Info("Lang." + ctx.I18n.Lang)
	}

	// fix meta link suffix
	for _, n := range ctx.Navs {
		n.Link = fixSuffix(n.Link)
	}
	ctx.Navs.I18n(ctx.I18n)

	// get owner, fix owner avatar link
	for _, a := range ctx.Authors {
		if a.IsOwner {
			ctx.Owner = a
		}
		a.Avatar = string(replacer([]byte(a.Avatar)))
	}
}

// read contents, including posts and pages
func (b *Builder) readContents(ctx *Context) {
	var (
		replacer     = replaceGlobalVars(b, ctx)
		htmlReplacer = replaceHTMLVars(b, ctx)
		filter       = func(p string) bool {
			return path.Ext(p) == ".md"
		}
	)
	postData, infoData, err := b.parseDir("post", filter)
	if err != nil {
		ctx.Error = err
		return
	}
	for k, blocks := range postData {
		post, err := model.NewPost(blocks, infoData[k])
		if err != nil {
			ctx.Error = err
			return
		}
		// if author name can find in ctx.Authors, use this named author,
		// if nil author but owner is set, use owner as author
		if post.Author != nil {
			if author, ok := ctx.Authors[post.Author.Name]; ok {
				post.Author = author
			}
		} else {
			if ctx.Owner != nil {
				post.Author = ctx.Owner
			}
		}
		post.Slug = string(replacer([]byte(post.Slug)))
		post.Thumb = string(replacer([]byte(post.Thumb)))
		post.PreviewHTML = template.HTML(htmlReplacer([]byte(post.PreviewHTML)))
		post.ContentHTML = template.HTML(htmlReplacer([]byte(post.ContentHTML)))
		ctx.Posts = append(ctx.Posts, post)
	}
	sort.Sort(model.Posts(ctx.Posts))

	ctx.Tags = make(map[string]*model.Tag)
	ctx.tagPosts = make(map[string][]*model.Post)
	for _, p := range ctx.Posts {
		for i, t := range p.Tags {
			ctx.Tags[t.Name] = p.Tags[i]
			ctx.tagPosts[t.Name] = append(ctx.tagPosts[t.Name], p)
		}
	}

	pageData, infoData, err := b.parseDir("page", filter)
	if err != nil {
		ctx.Error = err
		return
	}
	for k, blocks := range pageData {
		page, err := model.NewPage(blocks, infoData[k])
		if err != nil {
			ctx.Error = err
			return
		}
		// use named author
		if author, ok := ctx.Authors[page.Author.Name]; ok {
			page.Author = author
		}
		page.Slug = string(replacer([]byte(page.Slug)))
		page.Thumb = string(replacer([]byte(page.Thumb)))
		page.ContentHTML = template.HTML(htmlReplacer([]byte(page.ContentHTML)))
		ctx.Pages = append(ctx.Pages, page)
	}

	ctx.PageNodes = model.NewPageNodeGroup(ctx.Pages)
}

// parse bytes to blocks
func (b *Builder) parseBytes(data []byte) ([]parser.Block, error) {
	p := b.getParser(data[:32]) // read first 32 bytes to find current parser
	if p == nil {
		return nil, ErrParserMissing
	}
	blocks, err := p.Parse(data)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

// parse file to blocks
func (b *Builder) parseFile(file string) ([]parser.Block, error) {
	file = path.Join(b.opt.SrcDir, file)
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return b.parseBytes(fileData)
}

// parse files to blocks, return a map with file and blocks
func (b *Builder) parseFiles(files ...string) (map[string][]parser.Block, error) {
	data := make(map[string][]parser.Block)
	for _, file := range files {
		blocks, err := b.parseFile(file)
		if err != nil {
			return nil, err
		}
		data[file] = blocks
	}
	return data, nil
}

// parse files in directory with filter , return a map with file path and blocks
func (b *Builder) parseDir(dir string, filter func(string) bool) (map[string][]parser.Block, map[string]os.FileInfo, error) {
	dir = path.Join(b.opt.SrcDir, dir)
	data := make(map[string][]parser.Block)
	infoData := make(map[string]os.FileInfo)
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if filter != nil {
			if !filter(p) {
				return nil
			}
		}
		fileData, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		pa := b.getParser(fileData[:32]) // read first 32 bytes to find current parser
		if pa == nil {
			return ErrParserMissing
		}
		blocks, err := pa.Parse(fileData)
		if err != nil {
			return err
		}
		data[p] = blocks
		infoData[p] = info
		return nil
	})
	return data, infoData, err
}

// fix suffix to url,
// must append suffix
func fixSuffix(u string) string {
	// if url has host, full path
	if ur, _ := url.Parse(u); ur != nil {
		if ur.Host != "" {
			return u
		}
	}
	if u == "/" {
		return u
	}
	if path.Ext(u) == ".html" {
		return u
	}
	return u + ".html"
}

// global vars replacer
func replaceGlobalVars(b *Builder, ctx *Context) func([]byte) []byte {
	return func(str []byte) []byte {
		replacer := strings.NewReplacer(
			"@media/", "/"+ctx.mediaPath+"/",
			"@static/", "/"+ctx.staticPath+"/",
		)
		return []byte(replacer.Replace(string(str)))
	}
}

// global vars replacer
func replaceMarkdownVars(b *Builder, ctx *Context) func([]byte) []byte {
	return func(data []byte) []byte {
		data = bytes.Replace(data, []byte("(@media/"), []byte("(/"+ctx.mediaPath+"/"), -1)
		data = bytes.Replace(data, []byte("(@static/"), []byte("(/"+ctx.staticPath+"/"), -1)
		return data
	}
}

// global vars replacer in HTML
func replaceHTMLVars(b *Builder, ctx *Context) func([]byte) []byte {
	return func(data []byte) []byte {
		data = bytes.Replace(data, []byte(`href="@media/`), []byte(`href="/`+ctx.mediaPath+"/"), -1)
		data = bytes.Replace(data, []byte(`href="@static/`), []byte(`href="/`+ctx.staticPath+"/"), -1)
		data = bytes.Replace(data, []byte(`src="@media/`), []byte(`src="/`+ctx.mediaPath+"/"), -1)
		data = bytes.Replace(data, []byte(`src="@static/`), []byte(`src="/`+ctx.staticPath+"/"), -1)
		return data
	}
}
