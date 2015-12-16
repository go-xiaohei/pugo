package deploy

import (
	"errors"
	"path"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo-static/app/builder"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/ini.v1"
)

const (
	TYPE_GIT = "git"
)

var (
	_ DeployTask = new(GitTask)

	ErrGitNotRepo      = errors.New("destination directory is not a git repository")
	gitMessageReplacer = strings.NewReplacer("{now}", time.Now().Format(time.RFC3339))
)

type (
	// Git Deployment task
	GitTask struct {
		name string
		opt  *GitOption
	}
	// git options
	GitOption struct {
		Directory string `ini:"directory"` // if set, use this value. otherwise, use ctx.DstDir
		RepoUrl   string `ini:"repo_url"`  // remote repository url
		Branch    string `ini:"branch"`    // remote repository branch name
		User      string `ini:"user"`      // remote repository username, may need for http or https
		Password  string `ini:"password"`  // remote repository user password,may need for http or https
		Message   string `ini:"message"`   // commit message, only support {now} time string
	}
)

// New GitTask with name and ini.Section options
func (gt *GitTask) New(name string, section *ini.Section) (DeployTask, error) {
	// create a new GitTask
	var (
		g = &GitTask{
			name: name,
			opt: &GitOption{
				Message: "Site Updated at {now}",
			},
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

// GitTask's name
func (g *GitTask) Name() string {
	return g.name
}

// GitTask's type
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
		log15.Error("Deploy.Git.Error", "error", stderr)
		return err
	}
	log15.Debug("Deploy.[" + g.opt.RepoUrl + "].AddAll")

	// commit message
	message := gitMessageReplacer.Replace(opt.Message)
	if _, stderr, err := com.ExecCmdDir(
		ctx.DstDir, "git", []string{"commit", "-m", message}...); err != nil {
		log15.Error("Deploy.Git.Error", "error", stderr)
		return err
	}
	log15.Debug("Deploy.[" + g.opt.RepoUrl + "].Commit.'" + message + "'")

	// change remote url
	if _, stderr, err := com.ExecCmdDir(ctx.DstDir, "git", []string{
		"remote", "set-url", "origin", opt.remoteUrl(),
	}...); err != nil {
		log15.Error("Deploy.Git.Error", "error", stderr)
		return err
	}
	// push to repo
	if _, stderr, err := com.ExecCmdDir(ctx.DstDir, "git", []string{
		"push", "--force", "origin", opt.Branch}...); err != nil {
		log15.Error("Deploy.Git.Error", "error", stderr)
		if stderr != "" {
			return errors.New(stderr)
		}
		return err
	}
	log15.Debug("Deploy.[" + g.opt.RepoUrl + "].Push")
	return nil
}
