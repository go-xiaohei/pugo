package vars

import (
	"io/ioutil"
	"strings"
)

const (
	// Name is app name
	Name = "PuGo"
	// Desc is app usage
	Desc = "A Fast Static Site Generator"
	// Version is app version number
	Version = "0.10.0 (beta)"
)

var (
	// Commit is the building hash of commit
	Commit = ""
)

func init() {
	commitByte, _ := ioutil.ReadFile("commit")
	if len(commitByte) > 0 {
		Commit = strings.TrimSpace(string(commitByte))
	}
}
