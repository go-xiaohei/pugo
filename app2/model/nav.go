package model

import (
	"errors"
	"net/url"
)

type (
	// Nav is item of navigation
	Nav struct {
		Link        string `toml:"link"`
		Title       string `toml:"title"`
		OriginTitle string `toml:"-"`
		IsBlank     bool   `toml:"blank"`
		IconClass   string `toml:"icon"`
		HoverClass  string `toml:"hover"`
		I18n        string `toml:"i18n"`
		IsRemote    bool   `toml:"-"`
	}
	// NavGroup is group if items of navigation
	NavGroup []*Nav
)

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
