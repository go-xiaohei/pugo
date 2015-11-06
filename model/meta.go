package model

import (
	"errors"
	"gopkg.in/ini.v1"
	"pugo/parser"
)

var (
	ErrMetaBlockWrong = errors.New("meta-blocks-wrong")
)

type Meta struct {
	Title    string
	Subtitle string
	Keyword  string
	Desc     string
	Domain   string
}

func NewMeta(blocks []parser.Block) (*Meta, error) {
	if len(blocks) != 1 {
		return nil, ErrMetaBlockWrong
	}
	iniF, err := ini.Load(blocks[0].Bytes())
	if err != nil {
		return nil, err
	}
	section := iniF.Section("meta")
	meta := &Meta{
		Title:    section.Key("title").String(),
		Subtitle: section.Key("subtitle").String(),
		Keyword:  section.Key("keyword").String(),
		Desc:     section.Key("desc").String(),
		Domain:   section.Key("domain").String(),
	}
	return meta, nil
}
