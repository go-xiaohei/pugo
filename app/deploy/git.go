package deploy

import (
	"errors"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo-static/app/builder"
	"path"
)

const (
	TYPE_GIT = "git"
)

var (
	// _ DeployTask = new(GitTask)

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
		Branch  string // remote repository branch name
		Message string // commit message, only support {now} time string
	}
)

// New GitTask with name and ini.Section options
func (gt *GitTask) New(conf string) (*GitTask, error) {
	// create a new GitTask
	var (
		g = &GitTask{
			name: "git",
			opt: &GitOption{
				Message: "Site Updated at {now}",
			},
		}
	)
	return g, nil
}

// GitTask's name
func (g *GitTask) Name() string {
	return g.name
}

// Git deployment action
func (g *GitTask) Do(b *builder.Builder, ctx *builder.Context) error {
	gitDir := path.Join(ctx.DstDir, ".git")
	if !com.IsDir(gitDir) {
		return ErrGitNotRepo
	}
	return nil
	/*
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
	*/
}
