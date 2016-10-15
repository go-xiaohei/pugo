package model

// Comment save unique values for third-party comment systems
type Comment struct {
	Disqus  string `toml:"disqus" ini:"disqus"`
	Duoshuo string `toml:"duoshuo" ini:"duoshuo"`
}

// IsOK return the comment setting is valid
func (c *Comment) IsOK() bool {
	return c.Disqus != "" || c.Duoshuo != ""
}
