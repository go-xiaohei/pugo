package migrate

import (
	"fmt"
	"github.com/codegangsta/cli"
	rss "github.com/jteeuwen/go-pkg-rss"
	"gopkg.in/inconshreveable/log15.v2"
	"net/url"
	"os"
	"strings"
)

const (
	TypeRSS = "RSS"
)

var (
	_ Task = new(RSSTask)

	ErrRSSProtocolWrong = fmt.Errorf("Migrate RSS need protocol 'rss+http://' or 'rss+https://'")
)

type (
	RSSTask struct {
		opt *RSSOption
	}
	RSSOption struct {
		Source   string
		IsRemote bool
	}
)

func (rs *RSSTask) Is(conf string) bool {
	return strings.HasPrefix(conf, "rss+")
}

func (rs *RSSTask) New(ctx *cli.Context) (Task, error) {
	u, err := url.Parse(ctx.String("src"))
	if err != nil {
		return nil, err
	}

	opt := &RSSOption{}
	if len(u.Scheme) <= 4 {
		return nil, ErrRSSProtocolWrong
	}
	// get real schema, to get remote rss source
	u.Scheme = u.Scheme[4:]
	opt.IsRemote = true
	opt.Source = u.String()

	return &RSSTask{
		opt: opt,
	}, nil
}

func (rs *RSSTask) Type() string {
	return TypeRSS
}

func (rs *RSSTask) Dir() string {
	return "dir"
}

func (rs *RSSTask) Do() error {
	feed := rss.New(10, true, chanHandler, itemHandler)
	log15.Debug("RSS.Read." + rs.opt.Source)
	if err := feed.Fetch(rs.opt.Source, nil); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %s: %s\n", rs.opt.Source, err)
		return err
	}
	return nil
}

func chanHandler(feed *rss.Feed, newChannel []*rss.Channel) {
	fmt.Printf("%d new channel(s) in %s\n", len(newChannel), feed.Url)
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newItem []*rss.Item) {
	fmt.Printf("%d new item(s) in %s\n", len(newItem), feed.Url)
}
