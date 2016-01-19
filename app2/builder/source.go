package builder

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/app2/model"
	"gopkg.in/inconshreveable/log15.v2"
)

type (
	// Source include all sources data
	Source struct {
		Meta *model.Meta
	}
)

func readSource(ctx *Context) {
	var (
		srcDir = ""
	)
	// todo : clean code
	ctx.From = "dir:///../source"
	if srcDir, ctx.Err = toDir(ctx.From); ctx.Err != nil {
		return
	}
	if !com.IsDir(srcDir) {
		ctx.Err = fmt.Errorf("Directory '%s' is missing", srcDir)
		return
	}
	log15.Debug("Build|Source|%s", srcDir)

	// read meta
	// then read posts,
	// then read pages
	metaFile := filepath.Join(srcDir, "meta.toml")
	if !com.IsFile(metaFile) {
		ctx.Err = fmt.Errorf("Meta.toml is missing")
	}

	metaAll, err := readMeta(metaFile)
	if err != nil {
		ctx.Err = err
		return
	}

	fmt.Printf("%#v\n", metaAll)
}

func readMeta(file string) (*model.MetaAll, error) {
	meta := &model.MetaAll{}
	if _, err := toml.DecodeFile(file, meta); err != nil {
		return nil, err
	}
	return meta, nil
}

func toDir(urlString string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	if u.Scheme == "dir" || u.Scheme == "file" {
		return strings.Trim(u.Path, "/"), nil
	}
	return "", errors.New("Directory need schema dir:// or file ://")
}
