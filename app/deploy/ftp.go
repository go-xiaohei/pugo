package deploy

import (
	"errors"

	"github.com/go-xiaohei/pugo-static/app/builder"
	"github.com/jlaffaye/ftp"
	"gopkg.in/ini.v1"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	TYPE_FTP = "ftp"
)

var (
	_ DeployTask = new(FtpTask)
)

type (
	FtpTask struct {
		name string
		opt  *FtpOption
	}
	FtpOption struct {
		url      *url.URL
		Address  string `ini:"address"`
		User     string `ini:"user"`
		Password string `ini:"password"`
	}
)

func (fopt *FtpOption) isValid() error {
	if fopt.Address == "" || fopt.User == "" || fopt.Password == "" {
		return errors.New("deploy to ft need addres, username and password")
	}
	u, err := url.Parse(fopt.Address)
	if err != nil {
		return err
	}
	fopt.url = u
	return nil
}

func (ft *FtpTask) New(name string, section *ini.Section) (DeployTask, error) {
	var (
		f = &FtpTask{
			name: name,
			opt:  &FtpOption{},
		}
		err error
	)
	if err = section.MapTo(f.opt); err != nil {
		return nil, err
	}
	if err = f.IsValid(); err != nil {
		return nil, err
	}
	return f, nil
}

func (ft *FtpTask) Name() string {
	return ft.name
}

func (ft *FtpTask) IsValid() error {
	return ft.opt.isValid()
}

func (ft *FtpTask) Type() string {
	return TYPE_FTP
}

func (ft *FtpTask) Do(b *builder.Builder, ctx *builder.Context) error {
	client, err := ftp.DialTimeout(ft.opt.url.Host, time.Second*10)
	if err != nil {
		return err
	}
	defer client.Quit()
	if ft.opt.User != "" {
		if err = client.Login(ft.opt.User, ft.opt.Password); err != nil {
			return err
		}
	}

	// move to destination directory
	ftpDir := ft.opt.url.Path
	list, err := client.NameList(".")
	if err != nil {
		return err
	}
	isFound := false
	for _, name := range list {
		if name == strings.Trim(ftpDir, "/") {
			isFound = true
		}
	}
	if !isFound {
		if err = client.MakeDir(ftpDir); err != nil {
			return err
		}
	}
	if err = client.ChangeDir(ftpDir); err != nil {
		return err
	}

	var (
		rel, toFile string
	)
	if err = ctx.Diff.Walk(func(name string, entry *builder.DiffEntry) error {
		rel, _ = filepath.Rel(ctx.DstDir, name)
		toFile = filepath.ToSlash(filepath.Join(ftpDir, rel))

		// check file, if not exist, upload it
		_, err := client.Retr(toFile)
		if e, ok := err.(*textproto.Error); ok {
			if e.Code == ftp.StatusPageTypeUnknown {
				if err = uploadFile(client, toFile, name); err != nil {
					return err
				}
			}
		}

		if entry.Behavior == builder.DIFF_KEEP {
			return nil
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func uploadFile(c *ftp.ServerConn, path, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.Stor(path, f)
}
