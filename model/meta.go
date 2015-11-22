package model

import (
	"errors"

	"github.com/go-xiaohei/pugo-static/parser"
	"net/url"
)

var (
	ErrMetaBlockWrong = errors.New("meta-blocks-wrong")
)

// Meta contains basic info in site
type Meta struct {
	Title    string `ini:"title"`
	Subtitle string `ini:"subtitle"`
	Keyword  string `ini:"keyword"`
	Desc     string `ini:"desc"`
	Domain   string `ini:"domain"`
	Root     string `ini:"root"`
	Base     string `ini:"-"`
}

// blocks to Meta
func NewMeta(blocks []parser.Block) (*Meta, error) {
	if len(blocks) != 1 {
		return nil, ErrMetaBlockWrong
	}
	block, ok := blocks[0].(parser.MetaBlock)
	if !ok {
		return nil, ErrMetaBlockWrong
	}
	meta := new(Meta)
	if err := block.MapTo("meta", meta); err != nil {
		return nil, err
	}
	u, _ := url.Parse(meta.Root)
	meta.Base = u.Path
	if meta.Base == "/" {
		meta.Base = ""
	}
	return meta, nil
}
