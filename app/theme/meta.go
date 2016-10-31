package theme

import (
	"github.com/BurntSushi/toml"
	"github.com/go-xiaohei/pugo/app/model"
	"gopkg.in/ini.v1"
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

	Authors []*model.Author `toml:"author" ini:"-"`
	Refs    []*metaRef      `toml:"ref" ini:"-"`

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
	if t == model.FormatINI {
		meta := new(Meta)
		iniObj, err := ini.Load(data)
		if err != nil {
			return nil, err
		}
		if err = iniObj.Section("DEFAULT").MapTo(meta); err != nil {
			return nil, err
		}
		for _, authorKey := range iniObj.Section("author").KeysHash() {
			s := iniObj.Section("author." + authorKey)
			if len(s.Keys()) == 0 {
				continue
			}
			author := new(model.Author)
			if err = s.MapTo(author); err != nil {
				return nil, err
			}
			if author.Name != "" {
				meta.Authors = append(meta.Authors, author)
			}
		}
		for _, refKey := range iniObj.Section("ref").KeysHash() {
			s := iniObj.Section("ref." + refKey)
			if len(s.Keys()) == 0 {
				continue
			}
			ref := new(metaRef)
			if err = s.MapTo(ref); err != nil {
				return nil, err
			}
			if ref.Name != "" {
				meta.Refs = append(meta.Refs, ref)
			}
		}
		return meta, nil
	}
	return nil, nil
}
