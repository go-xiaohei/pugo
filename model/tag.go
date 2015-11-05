package model

import "net/url"

type Tag struct {
	Name string
	Url  string
}

func NewTag(name string) Tag {
	return Tag{
		Name: name,
		Url:  url.QueryEscape(name),
	}
}
