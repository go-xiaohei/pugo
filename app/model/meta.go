package model

import (
	"fmt"
	"net/url"
)

type (
	// Meta is meta info of website
	Meta struct {
		Title    string `toml:"title"`
		Subtitle string `toml:"subtitle"`
		Keyword  string `toml:"keyword"`
		Desc     string `toml:"desc"`
		Domain   string `toml:"domain"`
		Root     string `toml:"root"`
		Cover    string `toml:"cover"`
		Language string `toml:"lang"`
		Path     string `toml:"-"`
	}
	// MetaAll is all datat structs in meta.toml
	MetaAll struct {
		Meta        *Meta       `toml:"meta"`
		NavGroup    NavGroup    `toml:"nav"`
		AuthorGroup AuthorGroup `toml:"author"`
		Comment     *Comment    `toml:"comment"`
		Analytics   *Analytics  `toml:"analytics"`
	}
)

func (m *Meta) normalize() error {
	if m.Root == "" || m.Domain == "" || m.Title == "" {
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
