package def

type (
	// HexoMeta is definition of hexo meta block
	HexoMeta struct {
		Layout     string      `yaml:"layout"`
		Title      string      `yaml:"title"`
		Created    string      `yaml:"date"`
		Updated    string      `yaml:"updated"`
		Comments   bool        `yaml:"comments"` // not support in pugo
		Tags       interface{} `yaml:"tags"`
		Categories interface{} `yaml:"categories"` // not support in pugo
		Permalink  []string    `yaml:"permalink"`  // not support in pugo
	}

	// HexoConfig is definition of _config.yaml
	// partical, useful items
	HexoConfig struct {
		Title    string `yaml:"title"`       // support
		Subtitle string `yaml:"subtitle"`    // support
		Desc     string `yaml:"description"` // support
		Author   string `yaml:"author"`      // support
		Language string `yaml:"language"`
		Timezone string `yaml:"timezone"`

		URL              string `yaml:"url"`  // support
		Root             string `yaml:"root"` // support
		Permalink        string `yaml:"permalink"`
		PermalinkDefault string `yaml:"permalink_default"`

		SourceDir   string `yaml:"source_dir"`
		PublicDir   string `yaml:"public_dir"`
		TagDir      string `yaml:"tag_dir"`
		ArchiveDir  string `yaml:"archive_dir"`
		CategoryDir string `yaml:"category_dir"`
		CodeDir     string `yaml:"code_dir"`
		I18nDir     string `yaml:"i18n_dir"`
		SkipRender  string `yaml:"skip_render"`

		DateFormat string `yaml:"date_format"`
		TimeFormat string `yaml:"time_format"`
	}
)
