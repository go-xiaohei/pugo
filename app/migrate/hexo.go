package migrate

import (
	"bytes"
	"fmt"
	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"gopkg.in/inconshreveable/log15.v2"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	TypeHexo = "Hexo"
)

var (
	_ Task = new(HexoTask)

	ErrHexoSourceDirMissing = fmt.Errorf("Hexo source directory is not found")
)

type (
	HexoTask struct {
		opt    *HexoOption
		result map[string]*bytes.Buffer
		err    error
	}
	HexoOption struct {
		Dest   string
		Source string
	}
	hexoSources struct {
		mediaFile []string
		posts     []string
		pages     map[string]string
	}
)

func (ht *HexoTask) Is(conf string) bool {
	return strings.HasPrefix(conf, "hexo://")
}

func (ht *HexoTask) New(ctx *cli.Context) (Task, error) {
	u, err := url.Parse(ctx.String("src"))
	if err != nil {
		return nil, err
	}

	opt := &HexoOption{
		Dest: ctx.String("dest"),
	}
	opt.Source = strings.TrimSuffix(u.Path, "/")
	if opt.Source == "" {
		opt.Source = strings.TrimSuffix(u.Host, "/")
	}
	opt.Source = path.Join(opt.Source, "source")
	if !com.IsDir(opt.Source) {
		log15.Error("Hexo.Source.Dir.[" + opt.Source + "].Missing")
		return nil, ErrHexoSourceDirMissing
	}

	return &HexoTask{
		opt: opt,
	}, nil
}

func (ht *HexoTask) Type() string {
	return TypeHexo
}

func (ht *HexoTask) Do() (map[string]*bytes.Buffer, error) {
	var s *hexoSources
	if s, ht.err = ht.readSources(); ht.err != nil {
		return nil, ht.err
	}
	fmt.Println(s)
	return ht.result, ht.err
}

func (ht *HexoTask) readSources() (*hexoSources, error) {
	sources := &hexoSources{
		pages: make(map[string]string),
	}
	err := filepath.Walk(ht.opt.Source, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		println(p)
		return nil
	})
	return sources, err
}
