package models

// Metadata is metadata of the site
type Metadata struct {
	Title    string `toml:"title" ini:"title" json:"title,omitempty"`
	Subtitle string `toml:"subtitle" ini:"subtitle" json:"subtitle,omitempty"`
	Keyword  string `toml:"keyword" ini:"keyword" json:"keyword,omitempty"`
	Desc     string `toml:"desc" ini:"desc" json:"desc,omitempty"`
	Domain   string `toml:"domain" ini:"domain" json:"domain,omitempty"`
	Root     string `toml:"root" ini:"root" json:"root,omitempty"`
	Language string `toml:"lang" ini:"lang" json:"language,omitempty"`
}

// MetaFile is meta data of all meta file
type MetaFile struct {
	Meta *Metadata `toml:"meta" ini:"meta"`
}
