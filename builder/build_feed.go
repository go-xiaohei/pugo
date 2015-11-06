package builder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func (b *Builder) feed(ctx *context, r *Report) {
	baseUrl := "http://" + ctx.Meta.Domain
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	buf.WriteString(`<rss version="2.0">`)
	buf.WriteString("<channel>")
	buf.WriteString(fmt.Sprintf("<title>%s</title>", ctx.Meta.Title+" - "+ctx.Meta.Subtitle))
	buf.WriteString(fmt.Sprintf("<link>%s</link>", baseUrl))
	buf.WriteString(fmt.Sprintf("<description>%s</description>", ctx.Meta.Desc))
	buf.WriteString(fmt.Sprintf("<lastBuildDate>%s</lastBuildDate>", time.Now().Format(time.RFC1123Z)))
	for _, p := range ctx.Posts {
		buf.WriteString("<item>")
		buf.WriteString(fmt.Sprintf("<title>%s</title>", p.Title))
		buf.WriteString(fmt.Sprintf("<link>%s</link>", baseUrl+p.Url))
		// buf.WriteString(fmt.Sprintf("<comments>%s</comments>", fi.Comments))
		buf.WriteString(fmt.Sprintf("<pubDate>%s</pubDate>", p.Created.Raw.Format(time.RFC1123Z)))
		for _, c := range p.Tags {
			buf.WriteString(fmt.Sprintf("<category>%s</category>", c.Name))
		}
		buf.WriteString(fmt.Sprintf("<description><![CDATA[ %s ]]></description>", p.ContentHTML()))
		buf.WriteString("</item>")
	}
	buf.WriteString("</channel>")
	buf.WriteString("</rss>")

	dstFile := path.Join(ctx.DstDir, "feed.xml")
	os.MkdirAll(path.Dir(dstFile), os.ModePerm)
	r.Error = ioutil.WriteFile(dstFile, buf.Bytes(), os.ModePerm)
}
