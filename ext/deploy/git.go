package deploy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/builder"
	"gopkg.in/inconshreveable/log15.v2"
	"net/url"
	"time"
)

var (
	gitScheme       = "git://"
	gitCommitLayout = "PUGO BUILD UPDATE - {t}"
)

// Git is deployment of git repository
type Git struct {
	gitRepo    string
	gitMessage string
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
	u, err := url.Parse(dir)
	if err != nil {
		return nil, err
	}
	fmt.Println(u.Query())
	commitMessage := gitCommitLayout
	if m, ok := u.Query()["commit"]; ok && len(m) > 0 {
		commitMessage = m[0]
	}

	dir = u.Path
	if !com.IsDir(dir) {
		return nil, fmt.Errorf("git repository '%s' is missing", dir)
	}
	if !com.IsDir(filepath.Join(dir, ".git")) {
		return nil, fmt.Errorf("directory '%s' is not a git repository", dir)
	}
	ctx.Copied.CleanIgnoreFile = append(ctx.Copied.CleanIgnoreFile, "README.md", "LICENSE", "readme.me")
	ctx.To = "dir://" + dir // build to git repository
	log15.Debug("Deploy|Git|To|%s", ctx.To)
	return &Git{
		gitRepo:    dir,
		gitMessage: commitMessage,
	}, nil
}

// Action do git deploy action with built Context
func (g *Git) Action(ctx *builder.Context) error {
	// git add -A
	log15.Debug("Deploy|Git|git add -A")
	_, errOut, err := com.ExecCmdDir(g.gitRepo, "git", "add", "-A")
	if err = returnGetError(errOut, err); err != nil {
		return err
	}

	// git commit -m "message"
	log15.Debug("Deploy|Git|git commit -m")
	message := strings.Replace(g.gitMessage, "{t}", time.Now().Format(time.RFC1123), 1)
	_, errOut, err = com.ExecCmdDir(g.gitRepo, "git", "commit", "-m", message)
	if err = returnGetError(errOut, err); err != nil {
		return err
	}

	// git push
	log15.Debug("Deploy|Git|git push -f")
	_, errOut, err = com.ExecCmdDir(g.gitRepo, "git", "push", "-f")
	if err = returnGetError(errOut, err); err != nil {
		return err
	}
	return nil
}

func returnGetError(errOut string, err error) error {
	if errOut != "" && strings.Contains(errOut, "fatal:") {
		return fmt.Errorf(errOut)
	}
	if err != nil {
		return err
	}
	return nil
}
