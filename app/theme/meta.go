package theme

import (
	"github.com/BurntSushi/toml"
	"github.com/go-xiaohei/pugo/app/model"
)

// Meta is description of theme
type Meta struct {
	Name string   `toml:"name" ini:"name"`
	Repo string   `toml:"repo" ini:"repo"`
	URL  string   `toml:"url" ini:"url"`
	Date string   `toml:"date" ini:"date"`
	Desc string   `toml:"desc" ini:"desc"`
	Tags []string `toml:"tags" ini:"-"`

	MinVersion string `toml:"min_version" ini:"min_version"`

	Authors []*model.Author `toml:"authors" ini:"-"`
	Refs    []*metaRef      `toml:"refs" ini:"-"`

	License    string `toml:"license" ini:"license"`
	LicenseURL string `toml:"license_url" ini:"license_url"`
}

type metaRef struct {
	Name string `toml:"name" ini:"name"`
	URL  string `toml:"url" ini:"url"`
	Repo string `toml:"repo" ini:"repo"`
}

// NewMeta parse bytes to theme meta
func NewMeta(data []byte, t model.FormatType) (*Meta, error) {
	if t == model.FormatTOML {
		meta := new(Meta)
		if err := toml.Unmarshal(data, meta); err != nil {
			return nil, err
		}
		return meta, nil
	}
	return nil, nil
}
