package model

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
	}
	// MetaAll is all datat structs in meta.toml
	MetaAll struct {
		Meta *Meta `toml:"meta"`
	}
)
