package model

import (
	"errors"
	"github.com/go-xiaohei/pugo-static/parser"
)

var (
	ErrMetaBlockWrong = errors.New("meta-blocks-wrong")
)

type Meta struct {
	Title    string `ini:"title"`
	Subtitle string `ini:"subtitle"`
	Keyword  string `ini:"keyword"`
	Desc     string `ini:"desc"`
	Domain   string `ini:"domain"`
}

func NewMeta(blocks []parser.Block) (*Meta, error) {
	if len(blocks) != 1 {
		return nil, ErrMetaBlockWrong
	}
	block, ok := blocks[0].(parser.MetaBlock)
	if !ok {
		return nil, ErrMetaBlockWrong
	}
	meta := new(Meta)
	if err := block.Meta("meta", meta); err != nil {
		return nil, err
	}
	return meta, nil
}
