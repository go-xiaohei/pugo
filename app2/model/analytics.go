package model

// Analytics save unique values for web analytics service
type Analytics struct {
	Google string `toml:"google"`
	Baidu  string `toml:"baidu"`
}
