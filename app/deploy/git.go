package deploy

import (
	"errors"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app/builder"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	// TypeGit is git task string
	TypeGit = "git"
)

var (
	_ Task = new(GitTask)

	// ErrGitNotRepo shows the directory is not git repository
	ErrGitNotRepo = errors.New("destination directory is not a git repository")
	// ErrGitNoBranch shows that it cant't read repository's branch
	ErrGitNoBranch = errors.New("can not read git respository's branch")

	// default git commit message
	gitMessageReplacer = strings.NewReplacer("{now}", time.Now().Format(time.RFC3339))
)

type (
	// GitTask is  git deployment task
	GitTask struct {
		name      string
		opt       *GitOption
		directory string
	}
	// GitOption is git options
	GitOption struct {
		url     *url.URL
		Branch  string // remote repository branch name
		Message string // commit message, only support {now} time string
	}
)

// New returns new GitTask with name and ini.Section options
func (gt *GitTask) New(conf string) (Task, error) {
	// create a new GitTask
	g := &GitTask{
		name: "git",
		opt: &GitOption{
			Message: "Site Updated at {now}",
		},
	}

	// parse git repo directory
	u, err := url.Parse(conf)
	if err != nil {
		return nil, err
	}
	dir := u.Host
	if dir == "" {
		return nil, errors.New("git deploy conf need be git://git_repository_directory")
	}
	g.directory = dir

	// set commit message
	if commit := u.Query().Get("commit"); commit != "" {
		g.opt.Message = commit
	}
	g.opt.url = u
	return g, nil
}

// Type returns GitTask's name
func (gt *GitTask) Type() string {
	return TypeGit
}

// Dir returns GitTask's destination directory
func (gt *GitTask) Dir() string {
	return gt.directory
}

// Is checks GitTask
func (gt *GitTask) Is(conf string) bool {
	return strings.HasPrefix(conf, "git://")
}

// read repository's branch
func (gt *GitTask) readRepo(dest string) error {
	content, _, err := com.ExecCmdDir(dest, "git", []string{"branch"}...)
	if err != nil {
		return err
	}
	contentData := strings.Split(content, "\n")
	for _, cnt := range contentData {
		if strings.HasPrefix(cnt, "*") {
			cntData := strings.Split(cnt, " ")
			gt.opt.Branch = cntData[len(cntData)-1]
			return nil
		}
	}
	return nil
}

// Do executes git deploy action
func (gt *GitTask) Do(b *builder.Builder, ctx *builder.Context) error {
	gitDir := path.Join(ctx.DstDir, ".git")
	if !com.IsDir(gitDir) {
		return ErrGitNotRepo
	}
	var err error
	if err = gt.readRepo(ctx.DstDir); err != nil {
		return err
	}
	if gt.opt.Branch == "" {
		return ErrGitNoBranch
	}

	// add files
	if _, stderr, err := com.ExecCmdDir(ctx.DstDir, "git", []string{"add", "--all"}...); err != nil {
		log15.Error("Deploy.Git.Error", "error", stderr)
		return err
	}
	log15.Debug("Deploy.Git.[" + gt.opt.Branch + "].AddFiles")

	// commit message
	message := gitMessageReplacer.Replace(gt.opt.Message)
	if _, stderr, err := com.ExecCmdDir(ctx.DstDir, "git", []string{"commit", "-m", message}...); err != nil {
		log15.Error("Deploy.Git.Error", "error", stderr)
		return err
	}
	log15.Debug("Deploy.Git.[" + gt.opt.Branch + "].Commit.'" + message + "'")

	// push to repo
	_, stderr, err := com.ExecCmdDir(ctx.DstDir, "git", []string{
		"push", "--force", "origin", gt.opt.Branch}...)
	if err != nil {
		log15.Error("Deploy.Git.Error", "error", stderr)
		if stderr != "" {
			return errors.New(stderr)
		}
		return err
	}
	log15.Debug("Deploy.Git.[" + gt.opt.Branch + "].Push")
	return nil
}
