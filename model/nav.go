package model

import (
	"errors"
	"github.com/go-xiaohei/pugo-static/parser"
	"gopkg.in/ini.v1"
)

var (
	ErrNavBlockWrong = errors.New("nav-blocks-wrong")
)

type Nav struct {
	Link    string `ini:"link"`
	Title   string `ini:"title"`
	IsBlank bool   `ini:"blank"`

	IconClass  string `ini:"icon"`
	HoverClass string `ini:"hover"`
	I18n       string `ini:"i18n"`
	SubNav     []*Nav `ini:"-"`

	IsSeparator bool `ini:"-"`
	IsHover     bool `ini:"-"`
}

type Navs []*Nav

func (navs Navs) Hover(name string) {
	for _, n := range navs {
		if n.HoverClass == name {
			n.IsHover = true
		}
	}
}

func (navs Navs) Reset() {
	for _, n := range navs {
		n.IsHover = false
	}
}

func NewNavs(blocks []parser.Block) (Navs, error) {
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
	return Navs(navs), nil
}

func section2Nav(s *ini.Section) *Nav {
	nav := new(Nav)
	if err := s.MapTo(nav); err != nil {
		return nil
	}
	if nav.Link == "" {
		return nil
	}
	return nav
}
