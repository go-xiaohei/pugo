package deploy

import (
	"errors"
	"path"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo-static/app/builder"
)

const (
	TYPE_GIT = "git"
)

var (
	ErrGitNotRepo = errors.New("destination directory is not a git repository")

	// git message replacer
	gitMessageReplacer = strings.NewReplacer("{now}", time.Now().Format(time.RFC3339))
)

// git option
type GitOption struct {
	RepoUrl  string `ini:"repo_url"`
	Branch   string `ini:"branch"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	Message  string `ini:"message"`
}

// is valid option
func (g *GitOption) isValid() error {
	if g.RepoUrl == "" || g.Branch == "" || g.Message == "" {
		return errors.New("deploy to git need repo url, branch name and message")
	}
	return nil
}

// if set user,password and via http,
// add remote-url with {user}:{password}
func (g *GitOption) remoteUrl() string {
	if strings.HasPrefix(g.RepoUrl, "http://") {
		return strings.Replace(g.RepoUrl, "http://", "http://"+g.User+":"+g.Password+"@", 1)
	}
	if strings.HasPrefix(g.RepoUrl, "https://") {
		return strings.Replace(g.RepoUrl, "https://", "https://"+g.User+":"+g.Password+"@", 1)
	}
	return g.RepoUrl
}

// Git deployment action
func Git(opt GitOption, ctx *builder.Context) error {
	// check git repo
	gitDir := path.Join(ctx.DstDir, ".git")
	if !com.IsDir(gitDir) {
		return ErrGitNotRepo
	}
	// add files
	if _, _, err := com.ExecCmdDirBytes(
		ctx.DstDir,
		"git",
		[]string{"add", "--all"}...); err != nil {
		return err
	}
	// commit message
	if _, _, err := com.ExecCmdDirBytes(
		ctx.DstDir,
		"git",
		[]string{
			"commit",
			"-m",
			gitMessageReplacer.Replace(opt.Message)}...); err != nil {
		return err
	}
	// change remote url
	if _, _, err := com.ExecCmdDir(ctx.DstDir, "git", []string{
		"remote", "set-url", "origin", opt.remoteUrl(),
	}...); err != nil {
		return err
	}
	// push to repo
	if _, stderr, err := com.ExecCmdDir(ctx.DstDir, "git", []string{
		"push",
		"--force",
		"origin", opt.Branch,
	}...); err != nil {
		if stderr != "" {
			return errors.New(stderr)
		}
		return err
	}
	return nil
}
