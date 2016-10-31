package model

// Build is settings for builder in meta file
type Build struct {
	DisablePost bool `toml:"disable_post" ini:"disable_post"`
	DisablePage bool `toml:"disable_page" ini:"disable_page"`

	PostDir  string `toml:"post_dir" ini:"post_dir"`
	PageDir  string `toml:"page_dir" ini:"page_dir"`
	LangDir  string `toml:"lang_dir" ini:"lang_dir"`
	MediaDir string `toml:"media_dir" ini:"media_dir"`
}
