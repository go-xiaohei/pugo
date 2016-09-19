package deploy

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/urfave/cli"
	"gopkg.in/inconshreveable/log15.v2"
)

// Git is deployment of git repository
type Git struct {
	Repo    string
	Message string
	Local   string
	Branch  string
}

// Command return git deploy command
func (g *Git) Command() cli.Command {
	return cli.Command{
		Name:  "git",
		Usage: "deploy via git push",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "local", Value: "dest", Usage: "local website directory"},
			cli.StringFlag{Name: "repo", Usage: "local repository directory"},
			cli.StringFlag{Name: "message", Usage: "pushing commit message"},
			cli.StringFlag{Name: "branch", Value: "master", Usage: "the remote branch that git push to"},
		},
		Action: func(ctx *cli.Context) {
			g2, err := g.Create(ctx)
			if err != nil {
				log15.Error("Git|Fail|%s", err.Error())
				return
			}
			if err = g2.Do(); err != nil {
				log15.Error("Git|Fail|%s", err.Error())
				return
			}
			log15.Info("Git|Finish")
		},
	}
}

// String return git deployment typename
func (g *Git) String() string {
	return "Git"
}

// Create creates git deploy settings in Context
func (g *Git) Create(ctx *cli.Context) (Method, error) {
	g2 := &Git{
		Repo:    ctx.String("repo"),
		Message: ctx.String("message"),
		Local:   ctx.String("local"),
		Branch:  ctx.String("branch"),
	}
	if !com.IsDir(g2.Local) {
		return nil, fmt.Errorf("directory '%s' is not existed", g2.Local)
	}
	if !com.IsDir(g2.Repo) || !com.IsDir(filepath.Join(g2.Repo, ".git")) {
		return nil, fmt.Errorf("directory '%s' is not a git repository", g2.Repo)
	}
	if g2.Message == "" {
		g2.Message = "PUGO BUILD UPDATE - {t}"
	}
	return g2, nil
}

// Do do git deploy action with built Context
func (g *Git) Do() error {
	log15.Debug("Git|Overwrite")
	err := filepath.Walk(g.Local, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(g.Local, path)
		rel = filepath.ToSlash(rel)

		toFile := filepath.Join(g.Repo, rel)

		// overwrite file in repo
		os.MkdirAll(filepath.Dir(toFile), os.ModePerm)
		fileBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(toFile, fileBytes, os.ModePerm); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	log15.Debug("Git|git add -A")
	if err = runGitExec(g.Repo, []string{"add", "-A"}); err != nil {
		return err
	}
	log15.Debug("Git|git commit -m")
	message := strings.Replace(g.Message, "{t}", time.Now().Format(time.RFC1123), 1)
	if err = runGitExec(g.Repo, []string{"commit", "-m", message}); err != nil {
		return err
	}
	log15.Debug("Git|git push -f")
	if err = runGitExec(g.Repo, []string{"push", "-f"}); err != nil {
		return err
	}
	return nil
	/*
		// git add -A
		log15.Debug("Git|git add -A")
		_, errOut, err := com.ExecCmdDir(g.gitRepo, "git", "add", "-A")
		if err = returnGetError(errOut, err); err != nil {
			return err
		}

		// git commit -m "message"
		log15.Debug("Git|git commit -m")
		message := strings.Replace(g.gitMessage, "{t}", time.Now().Format(time.RFC1123), 1)
		_, errOut, err = com.ExecCmdDir(g.gitRepo, "git", "commit", "-m", message)
		if err = returnGetError(errOut, err); err != nil {
			return err
		}

		// git push
		log15.Debug("Git|git push -f")
		_, errOut, err = com.ExecCmdDir(g.gitRepo, "git", "push", "-f")
		if err = returnGetError(errOut, err); err != nil {
			return err
		}
		return nil
	*/
}

func runGitExec(dir string, commands []string) error {
	_, errOut, err := com.ExecCmdDir(dir, "git", commands...)
	return returnGetError(errOut, err)
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
