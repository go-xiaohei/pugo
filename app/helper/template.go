package helper

import "html/template"

// Str2HTML converts string to html
func Str2HTML(str string) template.HTML {
	return template.HTML(str)
}

// Bytes2HTML converts bytes to html
func Bytes2HTML(data []byte) template.HTML {
	return template.HTML(data)
}
