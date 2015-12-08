package deploy

import (
	"errors"

	"github.com/go-xiaohei/pugo-static/app/builder"
)

const (
	TYPE_FTP = "ftp"
)

var ()

type FtpOption struct {
	Address  string `ini:"address"`
	User     string `ini:"user"`
	Password string `ini:"password"`
}

func (fopt *FtpOption) isValid() error {
	if fopt.Address == "" || fopt.User == "" || fopt.Password == "" {
		return errors.New("deploy to ftp need addres, username and password")
	}
	return nil
}

func Ftp(opt FtpOption, ctx *builder.Context) error {
	return nil
}
