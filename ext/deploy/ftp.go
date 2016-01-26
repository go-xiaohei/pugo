package deploy

import (
	"net/url"
	"path"
	"strings"

	"github.com/go-xiaohei/pugo/app/builder"
	"github.com/goftp/ftp"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	ftpScheme = "ftp://"
)

// Ftp is ftp deployment
type Ftp struct {
	Address   string
	User      string
	Password  string
	Directory string
}

// Name is ftp deploy's name
func (f *Ftp) Name() string {
	return "FTP"
}

// Detect ftp deployment from Context
func (f *Ftp) Detect(ctx *builder.Context) (Task, error) {
	if !strings.HasPrefix(ctx.To, ftpScheme) {
		return nil, nil
	}
	u, err := url.Parse(ctx.To)
	if err != nil {
		return nil, err
	}

	f2 := &Ftp{
		Address: u.Host,
		User:    u.User.Username(),
	}
	f2.Password, _ = u.User.Password()
	f2.Directory = strings.Trim(u.Path, "/")
	log15.Debug("Deploy|FTP|To|%s", f2.Address)
	return f2, nil
}

// Action do ftp deploy process
func (f *Ftp) Action(ctx *builder.Context) error {
	return nil
}

// get dirs and subdirs
func getRecursiveDirs(dir string) []string {
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
