package model

import (
	"errors"
	"net/url"
	"path"

	"github.com/go-xiaohei/pugo/app/helper"
)

type (
	// Nav is item of navigation
	Nav struct {
		Link        string `toml:"link" ini:"link"`
		Title       string `toml:"title" ini:"title"`
		OriginTitle string `toml:"-" ini:"-"`
		IsBlank     bool   `toml:"blank" ini:"blank"`
		Icon        string `toml:"icon" ini:"icon"`
		IconClass   string `toml:"icon" ini:"icon"` // deprecated, old icon field name
		Hover       string `toml:"hover" ini:"hover"`
		I18n        string `toml:"i18n" ini:"i18n"`
		IsRemote    bool   `toml:"-" ini:"-"`
	}
	// NavGroup is group if items of navigation
	NavGroup []*Nav
)

// Tr print nav title with i18n helper
func (n *Nav) Tr(i18n *helper.I18n) string {
	return i18n.Tr("nav." + n.I18n)
}

// TrLink print nav link with i18n prefix
func (n *Nav) TrLink(i18n *helper.I18n) string {
	if n.IsRemote {
		return n.Link
	}
	return "/" + path.Join(i18n.Lang, n.Link)
}

// FixURL fix url path of all navigation items
func (ng NavGroup) FixURL(prefix string) {
	for _, n := range ng {
		if n.IsRemote {
			continue
		}
		n.Link = path.Join(prefix, n.Link)
	}
}

func (ng NavGroup) normalize() error {
	for _, n := range ng {
		if n.Link == "" || n.Title == "" {
			return errors.New("Nav's title or link is blank")
		}
		if u, _ := url.Parse(n.Link); u != nil && u.Host != "" {
			n.IsRemote = true
		}
		if n.I18n == "" {
			n.I18n = n.Title
		}
		n.OriginTitle = n.Title
	}
	return nil
}
