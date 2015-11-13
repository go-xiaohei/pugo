package command

import (
	"bytes"
	"fmt"
	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"gopkg.in/inconshreveable/log15.v2"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

var pathReplacer = strings.NewReplacer(".", "-", " ", "-", "+", "-")

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
				Name:  "post",
				Usage: "create new post",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "file",
						Usage: "filename of new post",
					},
				},
				Action: newPost(srcDir),
			},
			cli.Command{
				Name:   "page",
				Usage:  "create new page",
				Action: newPage(),
			},
		},
		HideHelp: true,
	}
}

func newSite() func(ctx *cli.Context) {
	return nil
}

func newPost(srcDir string) func(ctx *cli.Context) {
	return func(ctx *cli.Context) {
		file := ctx.String("file")
		if file == "" {
			log15.Crit("New Post need current filename, please run 'pugo new post --file=filename.md'")
		}
		file = path.Join(srcDir, "post", file)
		if com.IsFile(file) {
			log15.Crit("New Post file is existed..")
		}

		log15.Debug("NewPost." + file)
		var buf bytes.Buffer
		// ini block
		title := path.Base(file)
		fmt.Fprintln(&buf, "-----ini")
		fmt.Fprintf(&buf, `title = "%s"`+"\n", title)
		fmt.Fprintf(&buf, `slug = "%s"`+"\n", url.QueryEscape(pathReplacer.Replace(title)))
		fmt.Fprintf(&buf, `desc = "%s"`+"\n", title)
		fmt.Fprintf(&buf, "date = %s\n", time.Now().Format("2006-01-02 15:04"))
		fmt.Fprintf(&buf, "update_date = %s\n", time.Now().Format("2006-01-02 15:04"))
		fmt.Fprintln(&buf, "author = ")
		fmt.Fprintln(&buf, "author_email = ")
		fmt.Fprintln(&buf, "author_url = ")
		fmt.Fprintln(&buf, "tags = ")
		fmt.Fprintln(&buf, "")
		fmt.Fprintln(&buf, "-----markdown")
		fmt.Fprintln(&buf, "write your post content here")

		if err := ioutil.WriteFile(file, buf.Bytes(), os.ModePerm); err != nil {
			log15.Crit("NewPost.Fail", "error", err)
		}

		log15.Info("NewPost." + file + ".Success")
	}
}

func newPage() func(ctx *cli.Context) {
	return nil
}
