package helper

import (
	"bytes"
	"strings"

	"github.com/russross/blackfriday"
)

// copy from https://github.com/peachdocs/peach/blob/master/models/markdown.go
// thanks a lot
var (
	tab    = []byte("\t")
	spaces = []byte("    ")
)

// markdownRender sets some additions instead of default Render
type markdownRender struct {
	blackfriday.Renderer
}

// BlockCode overrides code block
func (mr *markdownRender) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	var tmp bytes.Buffer
	mr.Renderer.BlockCode(&tmp, text, strings.ToLower(lang))
	out.Write(bytes.Replace(tmp.Bytes(), tab, spaces, -1))
}

// Markdown converts markdown bytes to html bytes
func Markdown(raw []byte) []byte {
	htmlFlags := 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	renderer := &markdownRender{
		Renderer: blackfriday.HtmlRenderer(htmlFlags, "", ""),
	}

	extensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_AUTO_HEADER_IDS |
		blackfriday.EXTENSION_HEADER_IDS

	return blackfriday.Markdown(raw, renderer, extensions)
}

/*
func MarkdownBack(data []byte) ([]byte, error) {
	doc := htmldown.NewDocument(bytes.NewReader(data))
	str, err := doc.Markdown()
	return []byte(str), err
}
*/
