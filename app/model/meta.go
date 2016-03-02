package model

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/ini.v1"
)

type (
	// Meta is meta info of website
	Meta struct {
		Title    string `toml:"title" ini:"title"`
		Subtitle string `toml:"subtitle" ini:"subtitle"`
		Keyword  string `toml:"keyword" ini:"keyword"`
		Desc     string `toml:"desc" ini:"desc"`
		Domain   string `toml:"domain" ini:"domain"`
		Root     string `toml:"root" ini:"root"`
		Cover    string `toml:"cover" ini:"cover"`
		Language string `toml:"lang" ini:"lang"`
		Path     string `toml:"-" ini:"-"`
	}
	// MetaAll is all data struct in meta file
	MetaAll struct {
		Meta        *Meta       `toml:"meta"`
		NavGroup    NavGroup    `toml:"nav"`
		AuthorGroup AuthorGroup `toml:"author"`
		Comment     *Comment    `toml:"comment"`
		Analytics   *Analytics  `toml:"analytics"`
	}
)

// NewMetaAll parse bytes with correct FormatType
func NewMetaAll(data []byte, format FormatType) (*MetaAll, error) {
	var err error
	switch format {
	case FormatTOML:
		meta := &MetaAll{}
		if err = toml.Unmarshal(data, meta); err != nil {
			return nil, err
		}
		if err = meta.Normalize(); err != nil {
			return nil, err
		}
		return meta, nil
	case FormatINI:
		return newMetaAllFromINI(data)
	}
	return nil, fmt.Errorf("unsupported meta file format")
}

func newMetaAllFromINI(data []byte) (*MetaAll, error) {
	iniObj, err := ini.Load(data)
	if err != nil {
		return nil, err
	}
	section := iniObj.Section("meta")
	metaAll := new(MetaAll)

	// parse meta block
	meta := new(Meta)
	if err = section.MapTo(meta); err != nil {
		return nil, err
	}
	metaAll.Meta = meta

	// read navigation
	var (
		navGroup    []*Nav
		sectionName string
	)
	navKeys := iniObj.Section("nav").Keys()
	for _, k := range navKeys {
		sectionName = "nav." + k.Value()
		nav := new(Nav)
		if err = iniObj.Section(sectionName).MapTo(nav); err != nil {
			return nil, err
		}
		if nav.Title == "" && nav.Link == "" {
			continue
		}
		navGroup = append(navGroup, nav)
	}
	metaAll.NavGroup = navGroup

	// read author
	var authorGroup []*Author
	authorKeys := iniObj.Section("author").Keys()
	for _, k := range authorKeys {
		sectionName = "author." + k.Value()
		author := new(Author)
		if err = iniObj.Section(sectionName).MapTo(author); err != nil {
			return nil, err
		}
		if author.Name == "" {
			continue
		}
		authorGroup = append(authorGroup, author)
	}
	metaAll.AuthorGroup = authorGroup

	// read comment and analytics
	cmt := new(Comment)
	if err := iniObj.Section("comment").MapTo(cmt); err != nil {
		return nil, err
	}
	any := new(Analytics)
	if err := iniObj.Section("analytics").MapTo(cmt); err != nil {
		return nil, err
	}
	metaAll.Comment = cmt
	metaAll.Analytics = any

	if err = metaAll.Normalize(); err != nil {
		return nil, err
	}
	return metaAll, nil
}

// DomainURL return link with domain prefix
func (m *Meta) DomainURL(link string) string {
	link = strings.TrimPrefix(link, m.Path)
	link = strings.Trim(link, "/")
	return fmt.Sprintf("http://%s/%s", m.Domain, path.Join(strings.Trim(m.Path, "/"), link))
}

func (m *Meta) normalize() error {
	if (m.Root == "" && m.Domain == "") || m.Title == "" {
		return fmt.Errorf("meta title and (root or domain) cant be blank")
	}
	if m.Root == "" && m.Domain != "" {
		m.Root = "http://" + m.Domain + "/"
	}
	u, err := url.Parse(m.Root)
	if err != nil {
		return err
	}
	if m.Domain == "" {
		m.Domain = u.Host
	}
	m.Path = u.Path

	if m.Desc == "" {
		m.Desc = m.Title
	}
	if m.Keyword == "" {
		m.Keyword = m.Title
	}
	return nil
}

// Normalize make meta all data correct,
// it fills blank fields to correct values
func (ma *MetaAll) Normalize() error {
	var err error
	if err = ma.Meta.normalize(); err != nil {
		return err
	}
	if err = ma.NavGroup.normalize(); err != nil {
		return err
	}
	if err = ma.AuthorGroup.normalize(); err != nil {
		return err
	}
	return nil
}
