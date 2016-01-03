package deploy

import (
	"errors"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-xiaohei/pugo/app/builder"
	"github.com/goftp/ftp"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	// TypeFtp is FtpTask's type string
	TypeFtp = "ftp"
)

var (
	_ Task = new(FtpTask)
)

type (
	// FtpTask is ftp deployment task
	FtpTask struct {
		opt *FtpOption
	}
	// FtpOption is ftp deploy option
	FtpOption struct {
		url      *url.URL
		Address  string
		User     string
		Password string
	}
)

// New creates ftp task with section
func (ft *FtpTask) New(conf string) (Task, error) {
	u, err := url.Parse(conf)
	if err != nil {
		return nil, err
	}
	f := &FtpTask{opt: &FtpOption{Address: u.Host + u.Path}}
	if u.User != nil {
		f.opt.User = u.User.Username()
		f.opt.Password, _ = u.User.Password()
	}
	f.opt.url = u
	return f, nil
}

// Type returns ftp task's name
func (ft *FtpTask) Type() string {
	return TypeFtp
}

// Dir returns ftp task's real directory
func (ft *FtpTask) Dir() string {
	return path.Base(ft.opt.url.Path)
}

// Is checks ftp task
func (ft *FtpTask) Is(conf string) bool {
	return strings.HasPrefix(conf, "ftp://")
}

// Do executes ftp action
func (ft *FtpTask) Do(b *builder.Builder, ctx *builder.Context) error {
	client, err := ftp.DialTimeout(ft.opt.url.Host, time.Second*10)
	if err != nil {
		return err
	}
	log15.Debug("Deploy.[" + ft.opt.Address + "].Connect")
	defer client.Quit()
	if ft.opt.User != "" {
		if err = client.Login(ft.opt.User, ft.opt.Password); err != nil {
			return err
		}
	}
	log15.Debug("Deploy.[" + ft.opt.Address + "].Logged")
	ftpDir := strings.TrimPrefix(ft.opt.url.Path, "/")

	// change to UTF-8 mode
	if _, _, err = client.Exec(ftp.StatusCommandOK, "OPTS UTF8 ON"); err != nil {
		return err
	}

	if _, ok := client.Features()["UTF8"]; !ok {
		return errors.New("FTP server need utf-8 support")
	}

	// make dir
	makeFtpDir(client, getDirs(ftpDir))

	// change  to directory
	if err = client.ChangeDir(ftpDir); err != nil {
		return err
	}

	if b.Count < 2 {
		log15.Debug("Deploy.[" + ft.opt.Address + "].UploadAll")
		return ft.uploadAllFiles(client, ctx)
	}

	log15.Debug("Deploy.[" + ft.opt.Address + "].UploadDiff")
	return ft.uploadDiffFiles(client, ctx)
}

// upload files with checking diff status
func (ft *FtpTask) uploadDiffFiles(client *ftp.ServerConn, ctx *builder.Context) error {
	return ctx.Diff.Walk(func(name string, entry *builder.Entry) error {
		rel, _ := filepath.Rel(ctx.DstDir, name)
		rel = filepath.ToSlash(rel)

		if entry.Behavior == builder.DiffRemove {
			log15.Debug("Deploy.Ftp.Delete", "file", rel)
			return client.Delete(rel)
		}

		if entry.Behavior == builder.DiffKeep {
			if list, _ := client.List(rel); len(list) == 1 {
				// entry file should be older than uploaded file
				if entry.Time.Sub(list[0].Time).Seconds() < 0 {
					return nil
				}
			}
		}

		dirs := getDirs(path.Dir(rel))
		for i := len(dirs) - 1; i >= 0; i-- {
			client.MakeDir(dirs[i])
		}

		// upload file
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()
		if err = client.Stor(rel, f); err != nil {
			return err
		}
		log15.Debug("Deploy.Ftp.Stor", "file", rel)
		return nil
	})
}

// upload all files ignore diff status
func (ft *FtpTask) uploadAllFiles(client *ftp.ServerConn, ctx *builder.Context) error {
	var (
		createdDirs = make(map[string]bool)
		err         error
	)
	return ctx.Diff.Walk(func(name string, entry *builder.Entry) error {
		rel, _ := filepath.Rel(ctx.DstDir, name)
		rel = filepath.ToSlash(rel)

		// entry remove status, just remove it
		// the other files, just upload ignore diff status
		if entry.Behavior == builder.DiffRemove {
			log15.Debug("Deploy.Ftp.Delete", "file", rel)
			return client.Delete(rel)
		}

		// create directory recursive
		dirs := getDirs(path.Dir(rel))
		if len(dirs) > 0 {
			for i := len(dirs) - 1; i >= 0; i-- {
				dir := dirs[i]
				if !createdDirs[dir] {
					if err = client.MakeDir(dir); err == nil {
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
		if err = client.Stor(rel, f); err != nil {
			return err
		}
		log15.Debug("Deploy.Ftp.Stor", "file", rel)
		return nil
	})
}

// get dirs and subdirs
func getDirs(dir string) []string {
	if dir == "." || dir == "/" {
		return nil
	}
	dirs := []string{dir}
	for {
		dir = path.Dir(dir)
		if dir == "." || dir == "/" {
			break
		}
		dirs = append(dirs, dir)
	}
	return dirs
}

// make ftp directories,
// need make sub and parent directories
func makeFtpDir(client *ftp.ServerConn, dirs []string) error {
	for i := len(dirs) - 1; i >= 0; i-- {
		if err := client.MakeDir(dirs[i]); err != nil {
			return err
		}
	}
	return nil
}
