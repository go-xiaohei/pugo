package model

import (
	"errors"

	"github.com/go-xiaohei/pugo-static/parser"
	"net/url"
	"strings"
	"sync"
)

var (
	ErrMetaBlockWrong = errors.New("meta-blocks-wrong")

	navLock sync.Mutex
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
func NewAllMeta(blocks []parser.Block) (meta *Meta, navbar Navs, cmt *Comment, err error) {
	if len(blocks) != 1 {
		err = ErrMetaBlockWrong
		return
	}
	block, ok := blocks[0].(parser.MetaBlock)
	if !ok {
		err = ErrMetaBlockWrong
		return
	}

	// build meta
	meta = new(Meta)
	if err = block.MapTo("meta", meta); err != nil {
		return
	}
	if meta.Root == "" {
		meta.Root = "http://" + meta.Domain
	} else {
		if strings.HasSuffix(meta.Root, "/") {
			meta.Root = strings.TrimSuffix(meta.Root, "/")
		}
	}
	u, _ := url.Parse(meta.Root)
	meta.Base = u.Path
	if meta.Base == "/" {
		meta.Base = ""
	}

	// build nav
	navs := make([]*Nav, 0)
	navKeys := block.Keys("nav")
	for _, k := range navKeys {
		k = block.Item("nav", k)
		nav := new(Nav)
		if err := block.MapTo("nav."+k, nav); err != nil {
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
	navbar = Navs(navs)

	// build comment
	cmt = new(Comment)

	// disqus
	disqus := new(CommentDisqus)
	if err = block.MapTo("comment.disqus", disqus); err != nil {
		return
	}
	if disqus.Site != "" {
		cmt.Disqus = disqus
	}

	return
}

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

// Comment options
type Comment struct {
	Disqus *CommentDisqus `ini:"disqus"`
}

// Comment options of Disqus
type CommentDisqus struct {
	Site string `ini:"site"`
}

// Comment pasred third-party comments system,
// return as disqus,duoshuo, or empty string
func (c *Comment) String() string {
	using := []string{}
	if c.Disqus != nil {
		using = append(using, "disqus")
	}
	return strings.Join(using, ",")
}

// IsOK means is comment enabled,
// not empty settings
func (c *Comment) IsOK() bool {
	if c.Disqus != nil {
		return true
	}
	return false
}
