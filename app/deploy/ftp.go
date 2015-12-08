package deploy

import (
	"errors"

	"github.com/go-xiaohei/pugo-static/app/builder"
	"gopkg.in/ini.v1"
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
		Address  string `ini:"address"`
		User     string `ini:"user"`
		Password string `ini:"password"`
	}
)

func (fopt *FtpOption) isValid() error {
	if fopt.Address == "" || fopt.User == "" || fopt.Password == "" {
		return errors.New("deploy to ftp need addres, username and password")
	}
	return nil
}

func (ftp *FtpTask) New(name string, section *ini.Section) (DeployTask, error) {
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
	if err = ftp.IsValid(); err != nil {
		return nil, err
	}
	return f, nil
}

func (ftp *FtpTask) Name() string {
	return ftp.name
}

func (ftp *FtpTask) IsValid() error {
	return ftp.opt.isValid()
}

func (ftp *FtpTask) Type() string {
	return TYPE_FTP
}

func (ftp *FtpTask) Do(b *builder.Builder, ctx *builder.Context) error {
	return nil
}
