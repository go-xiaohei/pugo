package builder

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/go-xiaohei/pugo-static/model"
	"github.com/go-xiaohei/pugo-static/parser"
)

var (
	ErrParserMissing = errors.New("Parser-Unknown")
)

// parse data to context, if parsed all data to context for renders
func (b *Builder) ReadData(ctx *Context) {
	if b.readMeta(ctx); ctx.Error != nil {
		return
	}
	if b.readContents(ctx); ctx.Error != nil {
		return
	}
	// change dst dir if meta root is sub directory
	if ctx.Meta.Base != "" {
		ctx.DstDir = path.Join(ctx.DstDir, ctx.Meta.Base)
	}
	// load theme after data reading finished
	b.render.SetFunc("url", func(str ...string) string {
		return path.Join(append([]string{ctx.Meta.Base}, str...)...)
	})
	b.render.SetFunc("fullUrl", func(str ...string) string {
		return ctx.Meta.Root + path.Join(str...)
	})
	theme, err := b.render.Load(b.opt.Theme)
	if err != nil {
		ctx.Error = err
		return
	}

	ctx.Theme = theme
}

// read meta data, from meta.md,nav.md and comment.md
// they are global values.
func (b *Builder) readMeta(ctx *Context) {
	blocks, err := b.parseFile("meta.md")
	if err != nil {
		ctx.Error = err
		return
	}

	ctx.Meta, ctx.Navs, ctx.Authors, ctx.Comment, err = model.NewAllMeta(blocks)
	if err != nil {
		ctx.Error = err
		return
	}
	for _, n := range ctx.Navs {
		n.Link = fixSuffix(n.Link)
	}
}

// read contents, including posts and pages
func (b *Builder) readContents(ctx *Context) {
	filter := func(p string) bool {
		return path.Ext(p) == ".md"
	}
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
		// use named author
		if author, ok := ctx.Authors[post.Author.Name]; ok {
			post.Author = author
		}
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
		ctx.Pages = append(ctx.Pages, page)
	}
}

// parse file to blocks
func (b *Builder) parseFile(file string) ([]parser.Block, error) {
	file = path.Join(b.opt.SrcDir, file)
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	p := b.getParser(fileData[:32]) // read first 32 bytes to find current parser
	if p == nil {
		return nil, ErrParserMissing
	}
	blocks, err := p.Parse(fileData)
	if err != nil {
		return nil, err
	}
	return blocks, nil
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
	if u == "/" {
		return u
	}
	if path.Ext(u) == ".html" {
		return u
	}
	return u + ".html"
}
