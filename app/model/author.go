package model

import (
	"errors"

	"github.com/go-xiaohei/pugo/app/helper"
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
		IsOwner bool   // must be the first author
	}
	// AuthorGroup is collection of Authors
	AuthorGroup []*Author
)

func (ag AuthorGroup) normalize() error {
	if len(ag) == 0 {
		return errors.New("Must add an author")
	}
	ag[0].IsOwner = true
	for _, a := range ag {
		if a.Name == "" {
			return errors.New("author must have name")
		}
		if a.Avatar == "" && a.Email != "" {
			a.Avatar = helper.Gravatar(a.Email, 0)
		}
	}
	return nil
}
