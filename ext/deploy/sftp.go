package deploy

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gopkg.in/inconshreveable/log15.v2"
)

type Sftp struct {
	Host      string
	User      string
	Password  string
	Directory string
	Local     string

	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

// Command return sftp deploy command
func (s *Sftp) Command() cli.Command {
	return cli.Command{
		Name:  "sftp",
		Usage: "deploy via SSH account with SFTP",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "local", Value: "public", Usage: "local website directory"},
			cli.StringFlag{Name: "user", Usage: "ssh account name"},
			cli.StringFlag{Name: "password", Usage: "ssh account password"},
			cli.StringFlag{Name: "host", Usage: "ssh server address"},
			cli.StringFlag{Name: "directory", Usage: "sftp directory"},
		},
		Action: func(ctx *cli.Context) {
			s2, err := s.Create(ctx)
			if err != nil {
				log15.Error("Deploy|SFTP|Fail|%s", err.Error())
				return
			}
			if err = s2.Do(); err != nil {
				log15.Error("Deploy|SFTP|Fail|%s", err.Error())
				return
			}
			log15.Info("Deploy|SFTP|Finish")
		},
	}
}

func (s *Sftp) String() string {
	return "SFTP"
}

func (s *Sftp) Create(ctx *cli.Context) (Method, error) {
	s2 := &Sftp{
		Host:      ctx.String("host"),
		User:      ctx.String("user"),
		Password:  ctx.String("password"),
		Directory: ctx.String("directory"),
		Local:     ctx.String("local"),
	}
	if !com.IsDir(s2.Local) {
		return nil, fmt.Errorf("directory '%s' is empty", s2.Local)
	}
	if s2.Host == "" {
		return nil, fmt.Errorf("host is empty")
	}
	if s2.User == "" || s2.Password == "" {
		log15.Warn("Deploy|SFTP|No user or password")
	}
	if strings.HasPrefix(s2.Directory, "/~") {
		s2.Directory = strings.TrimPrefix(s2.Directory, "/~/")
	}
	return s2, nil
}

func (s *Sftp) Do() error {
	if err := s.connect(); err != nil {
		return err
	}
	defer s.sftpClient.Close()
	defer s.sshClient.Close()
	log15.Debug("Deploy|SFTP|%s|Connect", s.Host)
	makeSftpDir(s.sftpClient, getRecursiveDirs(s.Directory))
	log15.Debug("Deploy|SFTP|UploadAll")
	return s.UploadAll(s.Local)
}

func (s *Sftp) connect() error {
	conf := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
	}
	client, err := ssh.Dial("tcp", s.Host, conf)
	if err != nil {
		return err
	}
	s.sshClient = client
	s.sftpClient, err = sftp.NewClient(client)
	return err
}

// upload files without checking diff status
func (s *Sftp) UploadAll(local string) error {
	return filepath.Walk(local, func(p string, info os.FileInfo, err error) error {
		rel, _ := filepath.Rel(local, p)
		rel = filepath.ToSlash(rel)

		makeSftpDir(s.sftpClient, getRecursiveDirs(filepath.Dir(rel)))

		// upload file
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()
		f2, err := s.sftpClient.Create(path.Join(s.Directory, rel))
		if err != nil {
			return err
		}
		defer f2.Close()
		if _, err = io.Copy(f2, f); err != nil {
			return err
		}
		log15.Debug("Deploy|SFTO|Stor|%s", p)
		return nil
	})
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
