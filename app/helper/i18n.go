package helper

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/ini.v1"
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

// Trim trims string with lang prefix
func (i *I18n) Trim(str string) string {
	if strings.HasPrefix(str, "/"+i.Lang) {
		return strings.TrimLeft(str, "/"+i.Lang)
	}
	return strings.TrimLeft(str, "/")
}

// NewI18n reads toml bytes
func NewI18n(lang string, data []byte, ext string) (*I18n, error) {
	if ext == ".toml" {
		maps, err := I18nDataFromTOML(data)
		if err != nil {
			return nil, err
		}
		return &I18n{Lang: lang, values: maps}, nil
	}
	if ext == ".ini" {
		maps, err := I18nDataFromINI(data)
		if err != nil {
			return nil, err
		}
		return &I18n{Lang: lang, values: maps}, nil
	}
	return nil, nil
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

// I18nDataFromTOML parse toml data to map
func I18nDataFromTOML(data []byte) (map[string]map[string]string, error) {
	maps := make(map[string]map[string]string)
	if err := toml.Unmarshal(data, &maps); err != nil {
		return nil, err
	}
	for k, v := range maps {
		if len(v) == 0 {
			return nil, fmt.Errorf("i18n section '%s' is empty", k)
		}
	}
	return maps, nil
}

// I18nDataFromINI parse ini data to map
func I18nDataFromINI(data []byte) (map[string]map[string]string, error) {
	maps := make(map[string]map[string]string)
	iniObj, err := ini.Load(data)
	if err != nil {
		return nil, err
	}
	iniData := iniObj.Section("DEFAULT").KeysHash()
	for k, v := range iniData {
		k2 := strings.Split(k, ".")
		if len(k2) != 2 {
			continue
		}
		if _, ok := maps[k2[0]]; !ok {
			maps[k2[0]] = make(map[string]string)
		}
		maps[k2[0]][k2[1]] = v
	}
	for _, s := range iniObj.Sections() {
		if s.Name() == "DEFAULT" {
			continue
		}
		if m, ok := maps[s.Name()]; ok {
			for k, v := range s.KeysHash() {
				m[k] = v
			}
		} else {
			maps[s.Name()] = s.KeysHash()
		}
	}
	return maps, nil
}
