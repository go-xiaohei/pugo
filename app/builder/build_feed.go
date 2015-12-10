package builder

import (
	"bytes"
	"fmt"
	"github.com/Unknwon/com"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

// compile feed and sitemap
func (b *Builder) WriteFeed(ctx *Context) {
	baseUrl := strings.TrimSuffix(ctx.Meta.Root, "/")
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
	if ctx.Error = ioutil.WriteFile(dstFile, buf.Bytes(), os.ModePerm); ctx.Error != nil {
		return
	}
	if com.IsFile(dstFile) {
		ctx.Diff.Add(dstFile, DIFF_UPDATE, time.Now())
	} else {
		ctx.Diff.Add(dstFile, DIFF_ADD, time.Now())
	}

	// sitemap
	buf.Reset()
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	buf.WriteString(`<?xml-stylesheet type="text/xsl" href="` + ctx.Meta.Base + `/static/sitemap.xsl"?>`)
	buf.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	buf.WriteString("<url>")
	buf.WriteString(fmt.Sprintf("<loc>%s</loc>", baseUrl))
	buf.WriteString(fmt.Sprintf("<lastmod>%s/</lastmod>", time.Now().Format(time.RFC3339)))
	buf.WriteString("<changefreq>daily</changefreq>")
	buf.WriteString("<priority>1.0</priority>")
	buf.WriteString("</url>")

	for _, p := range ctx.Pages {
		buf.WriteString("<url>")
		buf.WriteString(fmt.Sprintf("<loc>%s</loc>", baseUrl+p.Url))
		buf.WriteString(fmt.Sprintf("<lastmod>%s/</lastmod>", p.Created.Raw.Format(time.RFC3339)))
		buf.WriteString("<changefreq>weekly</changefreq>")
		buf.WriteString("<priority>0.5</priority>")
		buf.WriteString("</url>")
	}

	for _, p := range ctx.Posts {
		buf.WriteString("<url>")
		buf.WriteString(fmt.Sprintf("<loc>%s</loc>", baseUrl+p.Url))
		buf.WriteString(fmt.Sprintf("<lastmod>%s/</lastmod>", p.Created.Raw.Format(time.RFC3339)))
		buf.WriteString("<changefreq>daily</changefreq>")
		buf.WriteString("<priority>0.6</priority>")
		buf.WriteString("</url>")
	}
	buf.WriteString("<url>")
	buf.WriteString(fmt.Sprintf("<loc>%s</loc>", baseUrl+"/archive.html"))
	buf.WriteString(fmt.Sprintf("<lastmod>%s/</lastmod>", time.Now().Format(time.RFC3339)))
	buf.WriteString("<changefreq>daily</changefreq>")
	buf.WriteString("<priority>0.6</priority>")
	buf.WriteString("</url>")

	for i := 1; i <= ctx.PostPageCount; i++ {
		buf.WriteString("<url>")
		buf.WriteString(fmt.Sprintf("<loc>%s/posts/%d.html</loc>", baseUrl, i))
		buf.WriteString(fmt.Sprintf("<lastmod>%s/</lastmod>", time.Now().Format(time.RFC3339)))
		buf.WriteString("<changefreq>daily</changefreq>")
		buf.WriteString("<priority>0.6</priority>")
		buf.WriteString("</url>")
	}

	for _, t := range ctx.Tags {
		buf.WriteString("<url>")
		buf.WriteString(fmt.Sprintf("<loc>%s</loc>", baseUrl+t.Url))
		buf.WriteString(fmt.Sprintf("<lastmod>%s/</lastmod>", time.Now().Format(time.RFC3339)))
		buf.WriteString("<changefreq>weekly</changefreq>")
		buf.WriteString("<priority>0.5</priority>")
		buf.WriteString("</url>")
	}

	buf.WriteString("</urlset>")
	dstFile = path.Join(ctx.DstDir, "sitemap.xml")
	os.MkdirAll(path.Dir(dstFile), os.ModePerm)
	ctx.Error = ioutil.WriteFile(dstFile, buf.Bytes(), os.ModePerm)
	if com.IsFile(dstFile) {
		ctx.Diff.Add(dstFile, DIFF_UPDATE, time.Now())
	} else {
		ctx.Diff.Add(dstFile, DIFF_ADD, time.Now())
	}
}
