package model

import (
	"errors"
	"net/url"
	"strings"
	"sync"

	"github.com/go-xiaohei/pugo/app/parser"
)

var (
	// ErrMetaBlockWrong is error of meta block
	ErrMetaBlockWrong = errors.New("meta-blocks-wrong")

	navLock sync.Mutex
)

type (
	// Meta contains basic info in site
	Meta struct {
		Title    string `ini:"title"`
		Subtitle string `ini:"subtitle"`
		Keyword  string `ini:"keyword"`
		Desc     string `ini:"desc"`
		Domain   string `ini:"domain"`
		Root     string `ini:"root"`
		Base     string `ini:"-"`
		Cover    string `ini:"cover"`
		Lang     string `ini:"lang"`
	}
	//MetaTotal contains all object in Meta
	MetaTotal struct {
		Meta    *Meta
		Nav     Navs
		Authors AuthorMap
		Comment *Comment
		Conf    *Conf
	}
)

// NewAllMeta parses blocks to MetaTotal
func NewAllMeta(blocks []parser.Block) (total MetaTotal, err error) {
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
	total = MetaTotal{}
	meta := new(Meta)
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
	if len(meta.Lang) == 0 {
		meta.Lang = "en" // use "en" as default
	}
	total.Meta = meta

	// build nav
	var navs []*Nav
	keys := block.Keys("nav")
	for _, k := range keys {
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
	total.Nav = Navs(navs)

	// build Authors
	authors := make(map[string]*Author)
	keys = block.Keys("author")
	for _, k := range keys {
		k = block.Item("author", k)
		author := new(Author)
		if err := block.MapTo("author."+k, author); err != nil {
			continue
		}
		if author.Name == "" {
			continue
		}
		if author.Nick == "" {
			author.Nick = author.Name
		}
		authors[k] = author
		if len(authors) == 1 {
			author.IsOwner = true
		}
	}
	total.Authors = AuthorMap(authors)

	// build comment
	cmt := new(Comment)

	// disqus
	if err = block.MapTo("comment", cmt); err != nil {
		return
	}
	if cmt.IsOK() {
		total.Comment = cmt
	}

	// conf
	cnf := new(Conf)
	hash := block.MapHash("build.ignore")
	for _, h := range hash {
		cnf.BuildIgnore = append(cnf.BuildIgnore, h)
	}
	total.Conf = cnf

	return
}

// Nav defines items in navigation
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

// Navs is collection of Nav
type Navs []*Nav

// Hover sets hover item
func (navs Navs) Hover(name string) {
	navLock.Lock()
	defer navLock.Unlock()
	for _, n := range navs {
		if n.HoverClass == name {
			n.IsHover = true
		}
	}
}

// Reset sets hover item to null
func (navs Navs) Reset() {
	navLock.Lock()
	defer navLock.Unlock()
	for _, n := range navs {
		n.IsHover = false
	}
}

// Comment options
type Comment struct {
	Disqus  string `ini:"disqus"`
	Duoshuo string `ini:"duoshuo"`
}

// Comment pasred third-party comments system,
// return as disqus,duoshuo, or empty string
func (c *Comment) String() string {
	using := []string{}
	if c.Disqus != "" {
		using = append(using, "disqus")
	}
	if c.Duoshuo != "" {
		using = append(using, "duoshuo")
	}
	return strings.Join(using, ",")
}

// IsOK means is comment enabled,
// not empty settings
func (c *Comment) IsOK() bool {
	if c.Disqus != "" || c.Duoshuo != "" {
		return true
	}
	return false
}

// Author of post or page
type Author struct {
	Name    string `ini:"name"`
	Nick    string `ini:"nick"`
	Email   string `ini:"email"`
	URL     string `ini:"url"`
	Avatar  string `ini:"avatar"` // todo: auto fill this field with gravatar
	IsOwner bool   // must be the first author
}

// AuthorMap is collection of Authors
type AuthorMap map[string]*Author

// Conf in meta, control building and deploying process
type Conf struct {
	BuildIgnore []string
}
