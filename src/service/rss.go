package service

import (
	"bytes"
	"github.com/go-xiaohei/pugo/src/model"
	"github.com/go-xiaohei/pugo/src/utils"
	"time"
)

var (
	RSS = new(RssService)
)

type RssService struct{}

func (rs *RssService) RSS(v interface{}) (*Result, error) {
	size := Setting.Content.RSSNumberLimit
	if size < 1 {
		size = 100
	}
	var (
		opt = ArticleListOption{
			Status:  model.ARTICLE_STATUS_PUBLISH,
			Order:   "create_time DESC",
			Page:    1,
			Size:    size,
			IsCount: false,
		}
		articles = make([]*model.Article, 0)
	)
	if err := Call(Article.List, opt, &articles); err != nil {
		return nil, err
	}

	var itemBuf bytes.Buffer
	for _, a := range articles {
		itemBuf.WriteString(`<item>
        <title>` + a.Title + `</title>
        <link>` + Setting.General.HostName + a.Href() + `</link>` + "\n")
		if Setting.Content.RSSFullText {
			itemBuf.WriteString(`<description><![CDATA[` + utils.Markdown2String(a.Body) + `]]></description>`)
		} else {
			itemBuf.WriteString(`<description><![CDATA[` + utils.Markdown2String(a.Preview) + `]]></description>`)
		}
		itemBuf.WriteString("\n" + `<pubDate>` + time.Unix(a.UpdateTime, 0).Format(time.RFC1123Z) + `</pubDate>
        <guid>` + Setting.General.HostName + a.Href() + `</guid>
    </item>`)
	}

	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
    <channel>
        <title>` + Setting.General.Title + `</title>
        <link>` + Setting.General.HostName + `</link>
        <description>` + Setting.General.FullTitle() + `</description>
        <pubDate>` + time.Unix(articles[0].UpdateTime, 0).Format(time.RFC1123Z) + `</pubDate>` + itemBuf.String() +
		`</channel>
</rss>`)
	return newResult(rs.RSS, &buf, &articles), nil
}
