package deploy

import (
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-xiaohei/pugo/app/builder"
	"github.com/go-xiaohei/pugo/app/model"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	sftpScheme = "sftp://"
)

type Sftp struct {
	Address   string
	User      string
	Password  string
	Directory string

	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

func (s *Sftp) Name() string {
	return "SFTP"
}

func (s *Sftp) Detect(ctx *builder.Context) (Task, error) {
	if !strings.HasPrefix(ctx.To, sftpScheme) {
		return nil, nil
	}
	u, err := url.Parse(ctx.To)
	if err != nil {
		return nil, err
	}
	s2 := &Sftp{
		Address: u.Host,
	}
	if u.User != nil {
		s2.User = u.User.Username()
		s2.Password, _ = u.User.Password()
	}
	s2.Directory = u.Path
	if strings.HasPrefix(s2.Directory, "/~") {
		s2.Directory = strings.TrimPrefix(s2.Directory, "/~/")
	}
	ctx.To = "dir://public"
	return s2, nil
}

func (s *Sftp) Action(ctx *builder.Context) error {
	if err := s.connect(); err != nil {
		return err
	}
	defer s.sftpClient.Close()
	defer s.sshClient.Close()
	log15.Debug("Deploy|SFTP|%s|Connect", s.Address)

	makeSftpDir(s.sftpClient, getRecursiveDirs(s.Directory))

	if builder.Counter() < 3 {
		log15.Debug("Deploy|SFTP|UploadAll")
		return s.UploadAll(ctx)
	}

	log15.Debug("Deploy|SFTP|UploadDiff")
	return s.UploadDiff(ctx)
}

func (s *Sftp) connect() error {
	conf := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
	}
	client, err := ssh.Dial("tcp", s.Address, conf)
	if err != nil {
		return err
	}
	s.sshClient = client
	s.sftpClient, err = sftp.NewClient(client)
	return err
}

// upload files without checking diff status
func (s *Sftp) UploadAll(ctx *builder.Context) error {
	for _, file := range ctx.Files.All() {
		rel, _ := filepath.Rel(ctx.DstDir(), file.URL)
		rel = filepath.ToSlash(rel)

		if file.Op == model.OpRemove {
			s.sftpClient.Remove(path.Join(s.Directory, rel))
			log15.Debug("Deploy|FTP|Remove|%s", file.URL)
			continue
		}

		makeSftpDir(s.sftpClient, getRecursiveDirs(filepath.Dir(rel)))

		// upload file
		f, err := os.Open(file.URL)
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
		log15.Debug("Deploy|SFTO|Stor|%s", file.URL)
	}
	return nil
}

// upload files with checking diff status
func (s *Sftp) UploadDiff(ctx *builder.Context) error {
	for _, file := range ctx.Files.All() {
		rel, _ := filepath.Rel(ctx.DstDir(), file.URL)
		rel = filepath.ToSlash(rel)

		if file.Op == model.OpKeep {
			log15.Debug("Deploy|SFTP|Skip|%s", file.URL)
			continue
		}
		if file.Op == model.OpRemove {
			s.sftpClient.Remove(path.Join(s.Directory, rel))
			log15.Debug("Deploy|FTP|Remove|%s", file.URL)
			continue
		}

		makeSftpDir(s.sftpClient, getRecursiveDirs(filepath.Dir(rel)))

		// upload file
		f, err := os.Open(file.URL)
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
		log15.Debug("Deploy|SFTO|Stor|%s", file.URL)
	}
	return nil
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
