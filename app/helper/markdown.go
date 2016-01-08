package helper

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/go-xiaohei/htmldown"
	"github.com/russross/blackfriday"
)

// copy from https://github.com/peachdocs/peach/blob/master/models/markdown.go
// thanks a lot
var (
	tab    = []byte("\t")
	spaces = []byte("    ")
)

// MarkdownRender sets some additions instead of default Render
type MarkdownRender struct {
	blackfriday.Renderer
}

// BlockCode overrides code block
func (mr *MarkdownRender) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	var tmp bytes.Buffer
	mr.Renderer.BlockCode(&tmp, text, strings.ToLower(lang))
	out.Write(bytes.Replace(tmp.Bytes(), tab, spaces, -1))
}

// Bytes2Markdown converts markdown bytes to html bytes
func Bytes2Markdown(raw []byte) []byte {
	htmlFlags := 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	renderer := &MarkdownRender{
		Renderer: blackfriday.HtmlRenderer(htmlFlags, "", ""),
	}

	extensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS

	return blackfriday.Markdown(raw, renderer, extensions)
}

// Bytes2MarkdownHTML converts markdown bytes to HTML
func Bytes2MarkdownHTML(raw []byte) template.HTML {
	return template.HTML(Bytes2Markdown(raw))
}

// HTML2Markdown converts to HTML string
func HTML2Markdown(html string) string {
	return htmldown.Markdown(html)
}
