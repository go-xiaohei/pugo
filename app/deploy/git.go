package deploy

import (
	"errors"
	"path"
	"strings"
	"time"

	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/ini.v1"

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

type (
	GitTask struct {
		name string
		opt  *GitOption
	}
	GitOption struct {
		Directory string `ini:"directory"`
		RepoUrl   string `ini:"repo_url"`
		Branch    string `ini:"branch"`
		User      string `ini:"user"`
		Password  string `ini:"password"`
		Message   string `ini:"message"`
	}
)

func (gt *GitTask) New(name string, section *ini.Section) (DeployTask, error) {
	// create a new GitTask
	var (
		g = &GitTask{
			name: name,
			opt:  &GitOption{},
		}
		err error
	)
	if err = section.MapTo(g.opt); err != nil {
		return nil, err
	}
	if err = g.IsValid(); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *GitTask) Name() string {
	return g.name
}

func (g *GitTask) Type() string {
	return TYPE_GIT
}

// is valid option
func (g *GitTask) IsValid() error {
	if g.opt.RepoUrl == "" || g.opt.Branch == "" || g.opt.Message == "" {
		return errors.New("deploy to git need repo url, branch name and message")
	}
	return nil
}

// if set user,password and via http,
// add remote-url with {user}:{password}
func (g *GitOption) remoteUrl() string {
	if g.User == "" || g.Password == "" {
		return g.RepoUrl
	}
	if strings.HasPrefix(g.RepoUrl, "http://") {
		return strings.Replace(g.RepoUrl, "http://", "http://"+g.User+":"+g.Password+"@", 1)
	}
	if strings.HasPrefix(g.RepoUrl, "https://") {
		return strings.Replace(g.RepoUrl, "https://", "https://"+g.User+":"+g.Password+"@", 1)
	}
	return g.RepoUrl
}

// Git deployment action
func (g *GitTask) Do(b *builder.Builder, ctx *builder.Context) error {
	opt := g.opt
	if opt.Directory == "" {
		opt.Directory = ctx.DstDir // use context destination directory as default
	}
	// check git repo
	gitDir := path.Join(opt.Directory, ".git")
	if !com.IsDir(gitDir) {
		return ErrGitNotRepo
	}
	// add files
	if _, stderr, err := com.ExecCmdDir(
		ctx.DstDir,
		"git",
		[]string{"add", "--all"}...); err != nil {
		log15.Debug("Deploy.Git.Error", "error", stderr)
		return err
	}
	// commit message
	if _, stderr, err := com.ExecCmdDir(
		ctx.DstDir, "git", []string{"commit", "-m",
			gitMessageReplacer.Replace(opt.Message)}...); err != nil {
		log15.Debug("Deploy.Git.Error", "error", stderr)
		return err
	}
	// change remote url
	if _, stderr, err := com.ExecCmdDir(ctx.DstDir, "git", []string{
		"remote", "set-url", "origin", opt.remoteUrl(),
	}...); err != nil {
		log15.Debug("Deploy.Git.Error", "error", stderr)
		return err
	}
	// push to repo
	if _, stderr, err := com.ExecCmdDir(ctx.DstDir, "git", []string{
		"push", "--force", "origin", opt.Branch}...); err != nil {
		log15.Debug("Deploy.Git.Error", "error", stderr)
		if stderr != "" {
			return errors.New(stderr)
		}
		return err
	}
	return nil
}
