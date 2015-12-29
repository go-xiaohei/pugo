package helper

import (
	"fmt"
	"gopkg.in/ini.v1"
)

type I18n struct {
	values map[string]string
}

func (i *I18n) Tr(str string) string {
	if v, ok := i.values[str]; ok {
		return v
	}
	return ""
}

func (i *I18n) Trf(str string, values ...interface{}) string {
	if v, ok := i.values[str]; ok {
		return fmt.Sprintf(v, values...)
	}
	return ""
}

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
	return &I18n{values: maps}, nil
}
