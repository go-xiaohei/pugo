package helper

import (
	"fmt"
	"html/template"
	"path"
	"strings"

	"gopkg.in/ini.v1"
)

// I18n object
type I18n struct {
	Lang   string // language string
	values map[string]string
}

// Tr converts string
func (i *I18n) Tr(str string) string {
	if v := i.values[str]; v != "" {
		return v
	}
	return str
}

// TrHTML converts string to html
func (i *I18n) TrHTML(str string) template.HTML {
	return template.HTML(i.Tr(str))
}

// Trf converts string with arguments
func (i *I18n) Trf(str string, values ...interface{}) string {
	if v, ok := i.values[str]; ok {
		return fmt.Sprintf(v, values...)
	}
	return str
}

// TrfHTML converts html string with arguments
func (i *I18n) TrfHTML(str string, values ...interface{}) template.HTML {
	return template.HTML(i.Trf(str, values...))
}

// NewI18n reads ini file with special key section
func NewI18n(file, key string) (*I18n, error) {
	f, err := ini.Load(file)
	if err != nil {
		return nil, err
	}
	if key == "" {
		key = "DEFAULT"
	}
	data := f.Section(key)
	maps := data.KeysHash()
	if len(maps) == 0 {
		return nil, nil
	}
	lang := path.Base(file)
	lang = strings.TrimSuffix(lang, path.Ext(lang))
	return &I18n{
		Lang:   lang,
		values: maps,
	}, nil
}

// NewI18nEmpty creates new empty i18n object,
// it will keep i18 tool working, but no translated value
func NewI18nEmpty() *I18n {
	return &I18n{
		Lang:   "nil",
		values: make(map[string]string),
	}
}

// NewI18nLanguageCode returns correct language code possibly
// en-US -> [en-US,en-us,en]
func NewI18nLanguageCode(lang string) []string {
	languages := []string{lang} // [en-US]
	lower := strings.ToLower(lang)
	if lower != lang {
		languages = append(languages, lower) // use lowercase language code, [en-us]
	}
	if strings.Contains(lang, "-") {
		languages = append(languages, strings.Split(lang, "-")[0]) // use first word if en-US, [en]
	}
	return languages
}
