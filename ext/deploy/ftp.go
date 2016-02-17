package deploy

import (
	"net/url"
	"path"
	"strings"

	"fmt"
	"github.com/go-xiaohei/pugo/app/builder"
	"github.com/go-xiaohei/pugo/app/model"
	"github.com/goftp/ftp"
	"gopkg.in/inconshreveable/log15.v2"
	"os"
	"path/filepath"
	"time"
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
	ctx.To = "dir://public" // reset
	return f2, nil
}

// Action do ftp deploy process
func (f *Ftp) Action(ctx *builder.Context) error {
	// connect to ftp
	client, err := ftp.DialTimeout(f.Address, time.Second*10)
	if err != nil {
		return err
	}
	log15.Debug("Deploy|FTP|%s|Connect", f.Address)
	defer client.Quit()
	if f.User != "" {
		if err = client.Login(f.User, f.Password); err != nil {
			return err
		}
	}

	// change to UTF-8 mode
	log15.Debug("Deploy|FTP|%s|UTF-8", f.Address)
	if _, _, err = client.Exec(ftp.StatusCommandOK, "OPTS UTF8 ON"); err != nil {
		return fmt.Errorf("OPTS UTF8 ON:%s", err.Error())
	}
	if _, ok := client.Features()["UTF8"]; !ok {
		return fmt.Errorf("FTP server need utf-8 support")
	}

	// make dir
	makeFtpDir(client, getRecursiveDirs(f.Directory))

	// change  to directory
	if err = client.ChangeDir(f.Directory); err != nil {
		return err
	}

	if builder.Counter() < 3 {
		log15.Debug("Deploy|FTP|UploadAll")
		return ftpUploadAll(client, ctx)
	}

	log15.Debug("Deploy|FTP|UploadDiff")
	return ftpUploadDiff(client, ctx)
}

// upload files without checking diff status
func ftpUploadAll(client *ftp.ServerConn, ctx *builder.Context) error {
	for _, file := range ctx.Files.All() {
		rel, _ := filepath.Rel(ctx.DstDir(), file.URL)
		rel = filepath.ToSlash(rel)

		if file.Op == model.OpRemove {
			client.Delete(rel)
			log15.Debug("Deploy|FTP|Remove|%s", file.URL)
			continue
		}

		makeFtpDir(client, getRecursiveDirs(filepath.Dir(rel)))

		// upload file
		f, err := os.Open(file.URL)
		if err != nil {
			return err
		}
		defer f.Close()
		if err = client.Stor(rel, f); err != nil {
			return err
		}
		log15.Debug("Deploy|FTP|Stor|%s", file.URL)
	}
	return nil
}

// upload files with checking diff status
func ftpUploadDiff(client *ftp.ServerConn, ctx *builder.Context) error {
	for _, file := range ctx.Files.All() {
		rel, _ := filepath.Rel(ctx.DstDir(), file.URL)
		rel = filepath.ToSlash(rel)

		if file.Op == model.OpKeep {
			log15.Debug("Deploy|FTP|Skip|%s", file.URL)
			continue
		}
		if file.Op == model.OpRemove {
			client.Delete(rel)
			log15.Debug("Deploy|FTP|Remove|%s", file.URL)
			continue
		}

		makeFtpDir(client, getRecursiveDirs(filepath.Dir(rel)))

		// upload file
		f, err := os.Open(file.URL)
		if err != nil {
			return err
		}
		defer f.Close()
		if err = client.Stor(rel, f); err != nil {
			return err
		}
		log15.Debug("Deploy|FTP|Stor|%s", file.URL)
	}
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
