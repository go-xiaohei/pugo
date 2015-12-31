package helper

import "html/template"

func Str2HTML(str string) template.HTML {
	return template.HTML(str)
}

func Bytes2HTML(data []byte) template.HTML {
	return template.HTML(data)
}
