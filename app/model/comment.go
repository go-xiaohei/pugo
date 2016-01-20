package model

// Comment save unique values for third-party comment systems
type Comment struct {
	Disqus  string `toml:"disqus"`
	Duoshuo string `toml:"duoshuo"`
}
