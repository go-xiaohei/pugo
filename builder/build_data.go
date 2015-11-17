package builder

import (
	"errors"
	"github.com/go-xiaohei/pugo-static/model"
	"github.com/go-xiaohei/pugo-static/parser"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
)

var (
	ErrParserMissing = errors.New("Parser-Unknown")
)

// parse data to context, if parsed all data to context for renders
func (b *Builder) ReadData(ctx *Context, r *Report) {
	if b.readMeta(ctx, r); r.Error != nil {
		return
	}
	if b.readContents(ctx, r); r.Error != nil {
		return
	}
}

// read meta data, from meta.md,nav.md and comment.md
// they are global values.
func (b *Builder) readMeta(ctx *Context, r *Report) {
	blocksMap, err := b.parseFiles("meta.md", "nav.md", "comment.md")
	if err != nil {
		r.Error = err
		return
	}

	meta, err := model.NewMeta(blocksMap["meta.md"])
	if err != nil {
		r.Error = err
		return
	}
	ctx.Meta = meta

	navs, err := model.NewNavs(blocksMap["nav.md"])
	if err != nil {
		r.Error = err
		return
	}
	ctx.Navs = navs

	if ctx.isSuffixed {
		for _, n := range ctx.Navs {
			n.Link = fixSuffix(n.Link)
		}
	}

	cmt, err := model.NewComment(blocksMap["comment.md"])
	if err != nil {
		r.Error = err
		return
	}
	ctx.Comment = cmt
}

// read contents, including posts and pages
func (b *Builder) readContents(ctx *Context, r *Report) {
	filter := func(p string) bool {
		return path.Ext(p) == ".md"
	}
	postData, infoData, err := b.parseDir("post", filter)
	if err != nil {
		r.Error = err
		return
	}
	for k, blocks := range postData {
		post, err := model.NewPost(blocks, infoData[k])
		if err != nil {
			r.Error = err
			return
		}
		if ctx.isSuffixed {
			post.Url = fixSuffix(post.Url)
		}
		ctx.Posts = append(ctx.Posts, post)
	}
	sort.Sort(model.Posts(ctx.Posts))

	ctx.Tags = make(map[string]*model.Tag)
	ctx.tagPosts = make(map[string][]*model.Post)
	for _, p := range ctx.Posts {
		for i, t := range p.Tags {
			ctx.Tags[t.Name] = &p.Tags[i]
			if ctx.isSuffixed {
				ctx.Tags[t.Name].Url = fixSuffix(ctx.Tags[t.Name].Url)
			}
			ctx.tagPosts[t.Name] = append(ctx.tagPosts[t.Name], p)
		}
	}

	pageData, infoData, err := b.parseDir("page", filter)
	if err != nil {
		r.Error = err
		return
	}
	for k, blocks := range pageData {
		page, err := model.NewPage(blocks, infoData[k])
		if err != nil {
			r.Error = err
			return
		}
		if ctx.isSuffixed {
			page.Url = fixSuffix(page.Url)
		}
		ctx.Pages = append(ctx.Pages, page)
	}
}

// parse file to blocks
func (b *Builder) parseFile(file string) ([]parser.Block, error) {
	file = path.Join(b.srcDir, file)
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
	dir = path.Join(b.srcDir, dir)
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
