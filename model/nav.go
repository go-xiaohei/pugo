package model

import (
	"errors"
	"gopkg.in/ini.v1"
	"pugo/parser"
)

var (
	ErrNavBlockWrong = errors.New("nav-blocks-wrong")
)

type Nav struct {
	Link    string
	Title   string
	IsBlank bool

	IconClass string
	I18n      string
	SubNav    []*Nav

	IsSeparator bool
}

func NewNavs(blocks []parser.Block) ([]*Nav, error) {
	if len(blocks) != 1 {
		return nil, ErrNavBlockWrong
	}
	iniF, err := ini.Load(blocks[0].Bytes())
	if err != nil {
		return nil, err
	}
	navSection := iniF.Section("nav")
	navKeys := navSection.Keys()
	navs := make([]*Nav, 0)
	for _, k := range navKeys {
		subSection := iniF.Section(k.String())
		nav := section2Nav(subSection)
		if nav == nil {
			continue
		}

		sub := subSection.Key("sub").Strings(",")
		if len(sub) > 0 {
			for _, s := range sub {
				if s == "-" {
					nav.SubNav = append(nav.SubNav, &Nav{IsSeparator: true})
					continue
				}
				n := section2Nav(iniF.Section(s))
				if n != nil {
					nav.SubNav = append(nav.SubNav, n)
				}
			}
		}
		navs = append(navs, nav)
	}
	return navs, nil
}

func section2Nav(s *ini.Section) *Nav {
	link := s.Key("link").String()
	if link == "" {
		return nil
	}
	nav := &Nav{
		Link:      link,
		Title:     s.Key("title").String(),
		IconClass: s.Key("icon").String(),
		I18n:      s.Key("i18n").String(),
	}
	nav.IsBlank, _ = s.Key("blank").Bool()
	return nav
}
