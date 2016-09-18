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
	"github.com/go-xiaohei/pugo/app/asset"
	"github.com/go-xiaohei/pugo/app/model"
	"github.com/urfave/cli"
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
			newOnlyDocFlag,
		},
		Before: Before,
		Action: newContent,
	}

	titleReplacer = strings.NewReplacer(" ", "-", ",", "-", ".", "-", "。", "-", "，", "-")
)

func newContent(ctx *cli.Context) error {
	if len(ctx.Args()) == 0 {
		log15.Error("need params\nusage:\n pugo new [post|page|site]")
		return nil
	}
	var err error
	switch ctx.Args()[0] {
	case "site":
		err = newSite(ctx.Bool("doc"))
	case "post":
		err = newPost(ctx.Args()[1:], ctx.String("to"))
	case "page":
		err = newPage(ctx.Args()[1:], ctx.String("to"))
	default:
		log15.Error("unknown params\nusage:\n pugo new [post|page|site]")
		return nil
	}
	if err != nil {
		log15.Crit("New|%s|%s", ctx.Args()[0], err.Error())
	}
	return nil
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
	log15.Debug("New|Post|%s", toFile)

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
	if post.Title == "" {
		post.Title = fileKey
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

func newPage(args []string, dstDir string) error {
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
	toFile := filepath.Join(dstDir, fmt.Sprintf("page/%d", time.Now().Year()), fileName)
	log15.Debug("New|Page|%s", toFile)

	if com.IsFile(toFile) {
		return errors.New("File Exist")
	}

	page := &model.Page{
		Title:    strings.Join(args, " "),
		Slug:     fileKey,
		Desc:     strings.Join(args, " "),
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Update:   time.Now().Format("2006-01-02 15:04:05"),
		NavHover: fileKey,
		Template: "page.html",
		Lang:     "",
		Sort:     0,
		Meta: map[string]interface{}{
			"file": fileName,
		},
	}
	if page.Title == "" {
		page.Title = fileKey
	}

	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err = encoder.Encode(page); err != nil {
		return err
	}

	buf2 := bytes.NewBufferString("```toml\n")
	buf2.Write(buf.Bytes())
	buf2.WriteString("```\n\n")
	buf2.WriteString("write you page content in " + fileName)

	os.MkdirAll(filepath.Dir(toFile), os.ModePerm)
	log15.Info("New|Page|Write|%s", toFile)
	return ioutil.WriteFile(toFile, buf2.Bytes(), os.ModePerm)
}

func newSite(onlyDoc bool) error {
	log15.Info("New|Extract|Assets")
	dirs := []string{"source", "theme", "doc"}
	isSuccess := true

	var (
		err       error
		isExtract = true
	)
	for _, dir := range dirs {
		isExtract = (dir != "doc")
		if onlyDoc {
			isExtract = (dir == "doc")
		}
		if !isExtract {
			continue
		}
		log15.Info("New|Extract|Directory|%s", dir)
		if err = asset.RestoreAssets("./", dir); err != nil {
			isSuccess = false
			break
		}
	}
	if !isSuccess {
		for _, dir := range dirs {
			os.RemoveAll(filepath.Join("./", dir))
		}
	}
	return err
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
