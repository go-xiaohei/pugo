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
		Meta        *Meta       `toml:"meta"`
		NavGroup    NavGroup    `toml:"nav"`
		AuthorGroup AuthorGroup `toml:"author"`
		Comment     *Comment    `toml:"comment"`
		Analytics   *Analytics  `toml:"analytics"`
	}
)

func (m *Meta) normalize() {
	if m.Root == "" && m.Domain != "" {
		m.Root = "http://" + m.Domain + "/"
	}
	if m.Desc == "" {
		m.Desc = m.Title
	}
	if m.Keyword == "" {
		m.Keyword = m.Title
	}
}

// Normalize make meta all data correct,
// it fills blank fields to correct values
func (ma *MetaAll) Normalize() error {
	var err error
	ma.Meta.normalize()
	if err = ma.NavGroup.normalize(); err != nil {
		return err
	}
	if err = ma.AuthorGroup.normalize(); err != nil {
		return err
	}
	return nil
}
