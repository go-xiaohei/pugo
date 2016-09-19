package deploy

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/goftp/ftp"
	"github.com/urfave/cli"
	"gopkg.in/inconshreveable/log15.v2"
)

// Ftp is ftp deployment
type Ftp struct {
	Local     string
	Host      string
	User      string
	Password  string
	Directory string
}

// Command return ftp deploy command
func (f *Ftp) Command() cli.Command {
	return cli.Command{
		Name:  "ftp",
		Usage: "deploy via FTP account",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "local", Value: "dest", Usage: "local website directory"},
			cli.StringFlag{Name: "user", Usage: "ftp account name"},
			cli.StringFlag{Name: "password", Usage: "ftp account password"},
			cli.StringFlag{Name: "host", Usage: "ftp host address"},
			cli.StringFlag{Name: "directory", Usage: "ftp directory"},
		},
		Action: func(ctx *cli.Context) {
			newFtp, err := f.Create(ctx)
			if err != nil {
				log15.Error("Ftp|Fail|%s", err.Error())
				return
			}
			if err = newFtp.Do(); err != nil {
				log15.Error("Ftp|Fail|%s", err.Error())
				return
			}
			log15.Info("Ftp|Finish")
		},
	}
}

// String is ftp deploy's name
func (f *Ftp) String() string {
	return "FTP"
}

// Create create ftp method from cli args
func (f *Ftp) Create(ctx *cli.Context) (Method, error) {
	ftpMethod := &Ftp{
		Local:     ctx.String("local"),
		Host:      ctx.String("host"),
		User:      ctx.String("user"),
		Password:  ctx.String("password"),
		Directory: ctx.String("directory"),
	}
	if !com.IsDir(ftpMethod.Local) {
		return nil, fmt.Errorf("%s is not directory", ftpMethod.Local)
	}
	if ftpMethod.Host == "" {
		return nil, fmt.Errorf("host is empty")
	}
	if ftpMethod.User == "" || ftpMethod.Password == "" {
		log15.Warn("Ftp|No user or password")
	}
	return ftpMethod, nil
}

// Do do ftp deploy process
func (f *Ftp) Do() error {
	// connect to ftp
	client, err := ftp.DialTimeout(f.Host, time.Second*10)
	if err != nil {
		return err
	}
	log15.Info("FTP|%s|Connect", f.Host)
	defer client.Quit()
	if f.User != "" {
		if err = client.Login(f.User, f.Password); err != nil {
			return err
		}
	}

	// change to UTF-8 mode
	log15.Debug("FTP|%s|UTF-8", f.Host)
	if _, _, err = client.Exec(ftp.StatusCommandOK, "OPTS UTF8 ON"); err != nil {
		if !strings.Contains(err.Error(), "No need to") { // sometimes show 202, no need to set UTF8 mode because always on
			return fmt.Errorf("OPTS UTF8 ON:%s", err.Error())
		}
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

	log15.Debug("FTP|UploadAll")
	return ftpUploadAll(client, f.Local)
}

// upload files without checking diff status
func ftpUploadAll(client *ftp.ServerConn, local string) error {
	return filepath.Walk(local, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(local, path)
		rel = filepath.ToSlash(rel)

		makeFtpDir(client, getRecursiveDirs(filepath.Dir(rel)))

		// upload file
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		if err = client.Stor(rel, f); err != nil {
			return err
		}
		log15.Debug("FTP|Stor|%s", path)
		return nil
	})
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
