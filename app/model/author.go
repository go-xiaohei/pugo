package model

import (
	"errors"

	"github.com/go-xiaohei/pugo/app/helper"
	"gopkg.in/ini.v1"
)

type (
	// Author is author item in meta file
	Author struct {
		Name    string `toml:"name" ini:"name"`
		Nick    string `toml:"nick" ini:"nick"`
		Email   string `toml:"email" ini:"email"`
		URL     string `toml:"url" ini:"url"`
		Avatar  string `toml:"avatar" ini:"avatar"`
		Bio     string `toml:"bio" ini:"bio"`
		Repo    string `toml:"repo" ini:"repo"` // github repository
		IsOwner bool   // must be the first author
	}
	// AuthorGroup is collection of Authors
	AuthorGroup []*Author
)

var (
	errAuthorInvalid    = errors.New("author must have name")
	errAuthorGroupEmpty = errors.New("must add an author")
)

func (a *Author) normalize() error {
	if a.Name == "" {
		return errAuthorInvalid
	}
	if a.Nick == "" {
		a.Nick = a.Name
	}
	if a.Avatar == "" && a.Email != "" {
		a.Avatar = helper.Gravatar(a.Email, 0)
	}
	return nil
}

func (ag AuthorGroup) normalize() error {
	if len(ag) == 0 {
		return errAuthorGroupEmpty
	}
	ag[0].IsOwner = true
	for _, a := range ag {
		if err := a.normalize(); err != nil {
			return err
		}
	}
	return nil
}

func newAuthorFromIniSection(section *ini.Section) (*Author, error) {
	a := &Author{
		Name:   section.Key("author").Value(),
		Email:  section.Key("author_email").Value(),
		URL:    section.Key("author_url").Value(),
		Avatar: section.Key("author_avatar").Value(),
		Bio:    section.Key("author_bio").Value(),
	}
	return a, a.normalize()
}
