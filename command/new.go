package command

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo-static/asset"
	"gopkg.in/inconshreveable/log15.v2"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

func New(srcDir, tplDir string) cli.Command {
	return cli.Command{
		Name:  "new",
		Usage: "create new site, post or page",
		Subcommands: []cli.Command{
			cli.Command{
				Name:   "site",
				Usage:  "create new site",
				Action: newSite(),
			},
			cli.Command{
				Name:   "post",
				Usage:  "create new post",
				Action: newPost(srcDir),
			},
			cli.Command{
				Name:   "page",
				Usage:  "create new page",
				Action: newPage(srcDir),
			},
		},
		HideHelp: true,
	}
}

func newSite() func(ctx *cli.Context) {
	return func(ctx *cli.Context) {
		log15.Info("New.Extract.Assets")
		dirs := []string{"source", "template"}
		isSuccess := true
		for _, dir := range dirs {
			if err := asset.RestoreAssets("./", dir); err != nil {
				log15.Error("New.Extract.Fail", "error", err)
				isSuccess = false
				break
			}
		}
		if !isSuccess {
			for _, dir := range dirs {
				os.RemoveAll(path.Join("./", dir))
			}
			return
		}
		log15.Info("New.Extract")
	}
}

func newPost(srcDir string) func(ctx *cli.Context) {
	return func(ctx *cli.Context) {
		file := time.Now().Format("2006-01-02-15-04") + ".md"
		file = path.Join(srcDir, "post", file)
		log15.Debug("New." + file + ".Begin")

		// write meta
		var buf bytes.Buffer
		// ini block
		title := strings.TrimSuffix(path.Base(file), path.Ext(file))
		fmt.Fprintln(&buf, "```ini")
		fmt.Fprintf(&buf, `title = "%s"`+"\n", title)
		fmt.Fprintf(&buf, `slug = "%s"`+"\n", url.QueryEscape(title))
		fmt.Fprintf(&buf, `desc = "%s"`+"\n", title)
		fmt.Fprintf(&buf, "date = %s\n", time.Now().Format("2006-01-02 15:04"))
		fmt.Fprintf(&buf, "update_date = %s\n", time.Now().Format("2006-01-02 15:04"))
		fmt.Fprintln(&buf, "author = author")
		fmt.Fprintln(&buf, "author_email = ")
		fmt.Fprintln(&buf, "author_url = ")
		fmt.Fprintln(&buf, "tags = post")
		fmt.Fprintln(&buf, "```")

		// write markdown content
		fmt.Fprintln(&buf, "")
		fmt.Fprintln(&buf, "write your post content here")

		// write to source file
		os.MkdirAll(path.Dir(file), os.ModePerm)
		if err := ioutil.WriteFile(file, buf.Bytes(), os.ModePerm); err != nil {
			log15.Crit("New.Fail", "error", err)
		}

		log15.Info("New." + file + "")
	}
}

func newPage(srcDir string) func(ctx *cli.Context) {
	return func(ctx *cli.Context) {
		file := time.Now().Format("2006-01-02-15-04") + ".md"
		file = path.Join(srcDir, "page", file)
		log15.Debug("New." + file + ".Begin")

		// write meta
		var buf bytes.Buffer
		// ini block
		title := strings.TrimSuffix(path.Base(file), path.Ext(file))
		log15.Debug("NewPage.Title." + title)
		fmt.Fprintln(&buf, "```ini")
		fmt.Fprintf(&buf, `title = "%s"`+"\n", title)
		fmt.Fprintf(&buf, `slug = "%s"`+"\n", url.QueryEscape(title))
		fmt.Fprintf(&buf, `desc = "%s"`+"\n", title)
		fmt.Fprintf(&buf, "date = %s\n", time.Now().Format("2006-01-02 15:04"))
		fmt.Fprintf(&buf, "update_date = %s\n", time.Now().Format("2006-01-02 15:04"))
		fmt.Fprintln(&buf, "author = author")
		fmt.Fprintln(&buf, "author_email = ")
		fmt.Fprintln(&buf, "author_url = ")
		fmt.Fprintln(&buf, "hover = ")
		fmt.Fprintln(&buf, "template = page.html")
		fmt.Fprintln(&buf, "```")

		// write markdown content
		fmt.Fprintln(&buf, "")
		fmt.Fprintln(&buf, "write your page content here")

		// write to source file
		os.MkdirAll(path.Dir(file), os.ModePerm)
		if err := ioutil.WriteFile(file, buf.Bytes(), os.ModePerm); err != nil {
			log15.Crit("New.Fail", "error", err)
		}

		log15.Info("New." + file)
	}
}
