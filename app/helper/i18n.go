package helper

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

// I18n object
type I18n struct {
	Lang   string // language string
	values map[string]map[string]string
}

// Tr converts string
func (i *I18n) Tr(str string) string {
	strSlice := strings.Split(str, ".")
	if len(strSlice) != 2 {
		return str
	}
	if m, ok := i.values[strSlice[0]]; ok {
		if v := m[strSlice[1]]; v != "" {
			return v
		}
	}
	return str
}

// Trf converts string with arguments
func (i *I18n) Trf(str string, values ...interface{}) string {
	return fmt.Sprintf(i.Tr(str), values...)
}

// NewI18n reads toml bytes
func NewI18n(lang string, data []byte) (*I18n, error) {
	maps := make(map[string]map[string]string)
	if err := toml.Unmarshal(data, &maps); err != nil {
		return nil, err
	}
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
		values: make(map[string]map[string]string),
	}
}

// LangCode returns correct language code possibly
// en-US -> [en-US,en-us,en]
func LangCode(lang string) []string {
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
