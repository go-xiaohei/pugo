package deploy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/builder"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	gitScheme = "git://"
)

// Git is deployment of git repository
type Git struct {
	gitRepo string
}

// Name return git deployment typename
func (g *Git) Name() string {
	return "Git"
}

// Detect detect git deploy settings in Context
func (g *Git) Detect(ctx *builder.Context) (Task, error) {
	if !strings.HasPrefix(ctx.To, gitScheme) {
		return nil, nil
	}
	dir := strings.TrimPrefix(ctx.To, gitScheme)
	if !com.IsDir(dir) {
		return nil, fmt.Errorf("git repository '%s' is missing", dir)
	}
	if !com.IsDir(filepath.Join(dir, ".git")) {
		return nil, fmt.Errorf("directory '%s' is not a git repository", dir)
	}
	ctx.To = "dir://public" // reset to public
	log15.Debug("Deploy|Git|To|%s", ctx.To)
	return &Git{
		gitRepo: dir,
	}, nil
}

// Action do git deploy action with built Context
func (g *Git) Action(ctx *builder.Context) error {
	return nil
}
