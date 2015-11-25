package model

import (
	"errors"
	"strings"

	"github.com/go-xiaohei/pugo-static/parser"
	"sync"
)

var (
	ErrNavBlockWrong = errors.New("nav-blocks-wrong")

	navLock sync.Mutex
)

// Nav defines items in navigatior
type Nav struct {
	Link    string `ini:"link"`
	Title   string `ini:"title"`
	IsBlank bool   `ini:"blank"`

	IconClass  string `ini:"icon"`
	HoverClass string `ini:"hover"`
	I18n       string `ini:"i18n"`

	SubNav      []*Nav `ini:"-"` // todo : no support yed
	IsSeparator bool   `ini:"-"`
	IsHover     bool   `ini:"-"`
}

type Navs []*Nav

func (navs Navs) Hover(name string) {
	navLock.Lock()
	defer navLock.Unlock()
	for _, n := range navs {
		if n.HoverClass == name {
			n.IsHover = true
		}
	}
}

func (navs Navs) Reset() {
	navLock.Lock()
	defer navLock.Unlock()
	for _, n := range navs {
		n.IsHover = false
	}
}

func NewNavs(blocks []parser.Block) (Navs, error) {
	if len(blocks) != 1 {
		return nil, ErrNavBlockWrong
	}
	block, ok := blocks[0].(parser.MetaBlock)
	if !ok {
		return nil, ErrMetaBlockWrong
	}
	navs := make([]*Nav, 0)
	navKeys := block.Keys("nav")
	for _, k := range navKeys {
		k = block.Item("nav", k)
		nav := new(Nav)
		if err := block.MapTo(k, nav); err != nil {
			continue
		}
		if nav.Link == "" {
			continue
		}
		sub := strings.Split(block.Item(k, "sub"), ",")
		if len(sub) > 0 && sub[0] != "" {
			for _, s := range sub {
				if s == "-" {
					nav.SubNav = append(nav.SubNav, &Nav{IsSeparator: true})
					continue
				}
				n2 := new(Nav)
				if err := block.MapTo(s, n2); err != nil {
					continue
				}
				if n2.Link == "" {
					continue
				}
				nav.SubNav = append(nav.SubNav, n2)
			}
		}
		navs = append(navs, nav)
	}
	return Navs(navs), nil
}
