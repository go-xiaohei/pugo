package vars

import (
	"io/ioutil"
	"strings"
)

var (
	// Commit is building hash of commit
	Commit = ""
)

func init() {
	if commitByte, _ := ioutil.ReadFile("commit"); len(commitByte) > 0 {
		Commit = strings.TrimSpace(string(commitByte))
	}
}

const (
	// Name is name of the application
	Name = "PuGo"
	// Desc is description of the application
	Desc = "A Simple Static Site Generator"
	// Version is version of the application
	Version = "0.11.0"
)
