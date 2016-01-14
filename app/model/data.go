package model

import "gopkg.in/ini.v1"

type (
	Data struct {
		keySlice   []string
		valueSlice []string
		valueMap   map[string]string
		children   map[string]*Data
	}
	DataGroup map[string]*Data
)

func NewData(file *ini.File) *Data {
	data := &Data{}
	for _, section := range file.Sections() {
		if section.Name() == "DEFAULT" {
			data.keySlice = section.KeyStrings()
			data.valueMap = section.KeysHash()
			continue
		}
		if len(data.children) == 0 {
			data.children = make(map[string]*Data)
		}
		d := &Data{
			keySlice: section.KeyStrings(),
			valueMap: section.KeysHash(),
		}
		data.children[section.Name()] = d
	}
	return data
}

func (d *Data) Value(key string) string {
	return d.valueMap[key]
}

func (d *Data) Keys() []string {
	return d.keySlice
}

func (d *Data) Values() []string {
	if len(d.valueSlice) != len(d.keySlice) {
		d.valueSlice = make([]string, len(d.keySlice))
		for i, _ := range d.valueSlice {
			d.valueSlice[i] = d.valueMap[d.keySlice[i]]
		}
	}
	return d.valueSlice
}

func (d *Data) Total() map[string]string {
	return d.valueMap
}

func (d *Data) Child(name string) *Data {
	if len(d.children) == 0 {
		return nil
	}
	return d.children[name]
}
