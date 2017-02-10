package deploy

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Unknwon/com"
	"github.com/pkg/sftp"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
	"gopkg.in/inconshreveable/log15.v2"
)

type Sftp struct {
	Host      string
	User      string
	Password  string
	Directory string
	Local     string

	KeySigner ssh.Signer

	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// Command return sftp deploy command
func (s *Sftp) Command() cli.Command {
	return cli.Command{
		Name:  "sftp",
		Usage: "deploy via SSH account with SFTP",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "local", Value: "dest", Usage: "local website directory"},
			cli.StringFlag{Name: "user", Usage: "ssh account name"},
			cli.StringFlag{Name: "password", Usage: "ssh account password"},
			cli.StringFlag{Name: "host", Usage: "ssh server address"},
			cli.StringFlag{Name: "pkey", Value: userHomeDir() + "/.ssh/id_rsa", Usage: "ssh private key"},
			cli.StringFlag{Name: "directory", Usage: "sftp directory"},
		},
		Action: func(ctx *cli.Context) {
			s2, err := s.Create(ctx)
			if err != nil {
				log15.Error("SFTP|Fail|%s", err.Error())
				return
			}
			if err = s2.Do(); err != nil {
				log15.Error("SFTP|Fail|%s", err.Error())
				return
			}
			log15.Info("SFTP|Finish")
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
		log15.Warn("SFTP|No user or password")
	}
	if s2.Directory == "" {
		s2.Directory = s2.Local
	}

	keyBytes, err := ioutil.ReadFile(ctx.String("pkey"))
	if err != nil {
		log15.Warn("SFTP|Fail to read ssh private key")
	}

	s2.KeySigner, err = ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		log15.Warn("SFTP|Fail to parse ssh private key")
	}

	if strings.HasPrefix(s2.Directory, "/~") {
		s2.Directory = strings.TrimPrefix(s2.Directory, "/~/")
	}
	return s2, nil
}

func (s *Sftp) Do() error {
	log15.Debug("SFTP|%s|Connect", s.Host)
	if err := s.connect(); err != nil {
		return err
	}
	defer s.sftpClient.Close()
	defer s.sshClient.Close()

	log15.Debug("SFTP|UploadAll%s", s.Local)
	return s.UploadAll(s.Local)
}

func (s *Sftp) connect() error {
	conf := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
			ssh.PublicKeys(s.KeySigner),
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

		remotePath := path.Join(s.Directory, rel)

		//check and create remote dir if needed
		if info.Mode().IsDir() {
			if _, err := s.sftpClient.Stat(remotePath); err == nil {
				return nil
			}

			if err := s.sftpClient.Mkdir(remotePath); err != nil {
				return fmt.Errorf("SFTP|Fail|create remote dir %s: %s", remotePath, err)
			}

			return nil
		}

		// upload file
		f, err := os.Open(p)
		if err != nil {
			return fmt.Errorf("SFTP|Fail|open remote file %s: %s", remotePath, err)
		}
		defer f.Close()
		f2, err := s.sftpClient.Create(remotePath)
		if err != nil {
			return fmt.Errorf("SFTP|Fail|create remote file %s: %s", remotePath, err)
		}
		defer f2.Close()
		if _, err = io.Copy(f2, f); err != nil {
			return fmt.Errorf("SFTP|Fail|upload file content %s: %s", remotePath, err)
		}
		log15.Debug("SFTP|Stor|%s", p)
		return nil
	})
}
