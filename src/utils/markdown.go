package utils

import (
	"github.com/russross/blackfriday"
	"html/template"
)

func Markdown2Bytes(str string) []byte {
	return blackfriday.MarkdownCommon([]byte(str))
}

func Markdown2String(str string) string {
	return string(blackfriday.MarkdownCommon([]byte(str)))
}

func Markdown2HTML(str string) template.HTML {
	return template.HTML(Markdown2String(str))
}
