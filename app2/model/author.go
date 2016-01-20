package model

import (
	"errors"

	"github.com/go-xiaohei/pugo/app2/helper"
)

type (
	// Author is author item in meta.toml
	Author struct {
		Name    string `toml:"name"`
		Nick    string `toml:"nick"`
		Email   string `toml:"email"`
		URL     string `toml:"url"`
		Avatar  string `toml:"avatar"`
		Bio     string `toml:"bio"`
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
