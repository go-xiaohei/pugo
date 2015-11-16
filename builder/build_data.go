package builder

import (
	"github.com/go-xiaohei/pugo-static/model"
	"github.com/go-xiaohei/pugo-static/parser"
	"os"
	"path"
	"path/filepath"
	"sort"
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
		ctx.Posts = append(ctx.Posts, post)
	}
	sort.Sort(model.Posts(ctx.Posts))

	ctx.Tags = make(map[string]*model.Tag)
	ctx.tagPosts = make(map[string][]*model.Post)
	for _, p := range ctx.Posts {
		for i, t := range p.Tags {
			ctx.Tags[t.Name] = &p.Tags[i]
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
		ctx.Pages = append(ctx.Pages, page)
	}
}

// parse file to blocks
func (b *Builder) parseFile(file string) ([]parser.Block, error) {
	file = path.Join(b.srcDir, file)
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	blocks, err := b.parser.ParseReader(f)
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
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		blocks, err := b.parser.ParseReader(f)
		if err != nil {
			return err
		}
		data[p] = blocks
		infoData[p] = info
		return nil
	})
	return data, infoData, err
}
