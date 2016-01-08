package deploy

import (
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-xiaohei/pugo/app/builder"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	// TypeSftp returns SftpTask's type
	TypeSftp = "sftp"
)

var (
	_ Task = new(SftpTask)
)

type (
	// SftpTask defines sftp deploy task
	SftpTask struct {
		opt *SftpOption
	}
	// SftpOption defines sftp options
	SftpOption struct {
		url       *url.URL
		Address   string
		User      string
		Password  string
		Directory string
	}
)

// New returns new sftp task with section
func (ft *SftpTask) New(conf string) (Task, error) {
	u, err := url.Parse(conf)
	if err != nil {
		return nil, err
	}
	f := &SftpTask{opt: &SftpOption{Address: u.Host + u.Path}}
	if u.User != nil {
		f.opt.User = u.User.Username()
		f.opt.Password, _ = u.User.Password()
	}
	f.opt.url = u

	p := f.opt.url.Path
	if strings.HasPrefix(p, "/~") {
		f.opt.Directory = strings.TrimPrefix(p, "/~/")
	} else {
		f.opt.Directory = p
	}
	return f, nil
}

// Type returns sftp task's name
func (ft *SftpTask) Type() string {
	return TypeSftp
}

// Is checks sftp task's name
func (ft *SftpTask) Is(conf string) bool {
	return strings.HasPrefix(conf, "sftp://")
}

// Dir returns sftp task's build directory
func (ft *SftpTask) Dir() string {
	return path.Base(ft.opt.Directory)
}

// Do executes sftp task do action
func (ft *SftpTask) Do(b *builder.Builder, ctx *builder.Context) error {
	conn, client, err := connectSftp(ft.opt)
	if err != nil {
		return err
	}
	defer client.Close()
	defer conn.Close()
	log15.Debug("Deploy.[" + ft.opt.Address + "].Connected")

	// just make directory, ignore error
	makeSftpDir(client, getDirs(ft.opt.Directory))

	if b.Count < 2 {
		log15.Debug("Deploy.[" + ft.opt.Address + "].UploadAll")
		return ft.uploadAllFiles(client, ctx)
	}

	log15.Debug("Deploy.[" + ft.opt.Address + "].UploadDiff")
	return ft.uploadDiffFiles(client, ctx)
}

// upload all files in context
func (ft *SftpTask) uploadAllFiles(client *sftp.Client, ctx *builder.Context) error {
	var (
		createdDirs = make(map[string]bool)
		err         error
	)
	return ctx.Diff.Walk(func(name string, entry *builder.Entry) error {
		rel, _ := filepath.Rel(ctx.DstDir, name)
		rel = filepath.ToSlash(rel)

		if entry.Behavior == builder.DiffRemove {
			log15.Debug("Deploy.Sftp.Delete", "file", rel)
			return client.Remove(path.Join(ft.opt.Directory, rel))
		}

		// create directory recursive
		dirs := getDirs(path.Dir(rel))
		if len(dirs) > 0 {
			for i := len(dirs) - 1; i >= 0; i-- {
				dir := dirs[i]
				if !createdDirs[dir] {
					if err = client.Mkdir(path.Join(ft.opt.Directory, dir)); err == nil {
						createdDirs[dir] = true
					}
				}
			}
		}

		// upload file
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()

		f2, err := client.Create(path.Join(ft.opt.Directory, rel))
		if err != nil {
			return err
		}
		defer f2.Close()

		if _, err = io.Copy(f2, f); err != nil {
			return err
		}
		log15.Debug("Deploy.Sftp.Stor", "file", rel)
		return nil
	})
}

// upload different context files
func (ft *SftpTask) uploadDiffFiles(client *sftp.Client, ctx *builder.Context) error {
	return ctx.Diff.Walk(func(name string, entry *builder.Entry) error {
		rel, _ := filepath.Rel(ctx.DstDir, name)
		rel = filepath.ToSlash(rel)

		if entry.Behavior == builder.DiffRemove {
			log15.Debug("Deploy.Sftp.Delete", "file", rel)
			return client.Remove(path.Join(ft.opt.Directory, rel))
		}

		target := path.Join(ft.opt.Directory, rel)
		if entry.Behavior == builder.DiffKeep {
			if fi, _ := client.Stat(target); fi != nil {
				// entry file should be older than uploaded file
				if entry.Time.Sub(fi.ModTime()).Seconds() < 0 {
					return nil
				}
			}
		}

		dirs := getDirs(path.Dir(rel))
		for i := len(dirs) - 1; i >= 0; i-- {
			client.Mkdir(path.Join(ft.opt.Directory, dirs[i]))
		}

		// upload file
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()

		f2, err := client.Create(target)
		if err != nil {
			return err
		}
		defer f2.Close()

		if _, err = io.Copy(f2, f); err != nil {
			return err
		}
		log15.Debug("Deploy.Sftp.Stor", "file", rel)
		return nil
	})
}

// connect to sftp, get ssh connection and sftp client
func connectSftp(opt *SftpOption) (*ssh.Client, *sftp.Client, error) {
	conf := &ssh.ClientConfig{
		User: opt.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(opt.Password),
		},
	}
	client, err := ssh.Dial("tcp", opt.url.Host, conf)
	if err != nil {
		return nil, nil, err
	}
	s, err := sftp.NewClient(client)
	return client, s, err
}

// make sftp directories
func makeSftpDir(client *sftp.Client, dirs []string) error {
	for i := len(dirs) - 1; i >= 0; i-- {
		if err := client.Mkdir(dirs[i]); err != nil {
			return err
		}
	}
	return nil
}
