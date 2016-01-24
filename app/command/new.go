package command

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/model"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	// New is command of 'new'
	New = cli.Command{
		Name:  "new",
		Usage: "create new content",
		Flags: []cli.Flag{
			newToFlag,
			debugFlag,
		},
		Before: Before,
		Action: newContent,
	}

	titleReplacer = strings.NewReplacer(" ", "-", ",", "-", ".", "-", "。", "-", "，", "-")
)

func newContent(ctx *cli.Context) {
	if len(ctx.Args()) == 0 {
		log15.Error("need params\nusage:\n pugo new [post|page|site]")
		return
	}
	var err error
	switch ctx.Args()[0] {
	case "post":
		err = newPost(ctx.Args()[1:], ctx.String("to"))
	default:
		log15.Error("unknown params\nusage:\n pugo new [post|page|site]")
		return
	}
	if err != nil {
		log15.Crit("New|%s|%s", ctx.Args()[0], err.Error())
	}
}

func newPost(args []string, dstDir string) error {
	dstDir, err := toDir(dstDir)
	if err != nil {
		return err
	}
	fileKey := time.Now().Format("01-02-15-04-05")
	if len(args) > 0 {
		fileKey = strings.Join(args, "-")
		fileKey = titleReplacer.Replace(fileKey)
	}
	fileName := fileKey + ".md"
	toFile := filepath.Join(dstDir, fmt.Sprintf("post/%d", time.Now().Year()), fileName)
	log15.Debug("New|Post|To|%s", toFile)

	if com.IsFile(toFile) {
		return errors.New("File Exist")
	}

	post := &model.Post{
		Title:     strings.Join(args, " "),
		Slug:      fileKey,
		Desc:      strings.Join(args, " "),
		Date:      time.Now().Format("2006-01-02 15:04:05"),
		Update:    time.Now().Format("2006-01-02 15:04:05"),
		TagString: []string{"tag"},
	}

	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err = encoder.Encode(post); err != nil {
		return err
	}

	buf2 := bytes.NewBufferString("```toml\n")
	buf2.Write(buf.Bytes())
	buf2.WriteString("```\n\n")
	buf2.WriteString("write you post content in " + fileName)

	os.MkdirAll(filepath.Dir(toFile), os.ModePerm)
	log15.Info("New|Post|Write|%s", toFile)
	return ioutil.WriteFile(toFile, buf2.Bytes(), os.ModePerm)
}

func toDir(urlString string) (string, error) {
	if !strings.Contains(urlString, "://") {
		return urlString, nil
	}
	if strings.HasPrefix(urlString, "dir://") {
		return strings.TrimPrefix(urlString, "dir://"), nil
	}
	return "", errors.New("Directory need schema dir://")
}
