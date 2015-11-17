package parser

import (
	"bytes"
	"gopkg.in/ini.v1"
)

const (
	BLOCK_INI = "ini"
)

type (
	MetaBlock interface {
		Block
		MapTo(section string, v interface{}) error
		MapHash(section string) map[string]string
		Keys(section string) []string
		Item(k1 string, keys ...string) string
	}
)

type IniBlock struct {
	data      []byte
	iniObject *ini.File
}

func (ib *IniBlock) New() Block {
	return new(IniBlock)
}

func (ib *IniBlock) Type() string {
	return BLOCK_INI
}

func (ib *IniBlock) Is(mark []byte) bool {
	return bytes.Equal(mark, []byte(BLOCK_INI))
}

func (ib *IniBlock) Write(data []byte) error {
	ib.data = append(ib.data, data...)
	ib.iniObject = nil // todo: need locker
	return nil
}

func (ib *IniBlock) Bytes() []byte {
	return bytes.TrimRight(ib.data, "\n")
}

func (ib *IniBlock) MapTo(section string, v interface{}) error {
	if ib.iniObject == nil {
		var err error
		ib.iniObject, err = ini.Load(ib.Bytes())
		if err != nil {
			return err
		}
	}
	if section == "" {
		section = "DEFAULT"
	}
	return ib.iniObject.Section(section).MapTo(v)
}

func (ib *IniBlock) MapHash(section string) map[string]string {
	if ib.iniObject == nil {
		var err error
		ib.iniObject, err = ini.Load(ib.Bytes())
		if err != nil {
			return map[string]string{}
		}
	}
	if section == "" {
		section = "DEFAULT"
	}
	return ib.iniObject.Section(section).KeysHash()
}

func (ib *IniBlock) Keys(section string) []string {
	if ib.iniObject == nil {
		var err error
		ib.iniObject, err = ini.Load(ib.Bytes())
		if err != nil {
			return []string{}
		}
	}
	if section == "" {
		section = "DEFAULT"
	}
	return ib.iniObject.Section(section).KeyStrings()
}

func (ib *IniBlock) Item(k1 string, keys ...string) string {
	if ib.iniObject == nil {
		var err error
		ib.iniObject, err = ini.Load(ib.Bytes())
		if err != nil {
			return ""
		}
	}
	if len(keys) == 0 {
		section := ib.iniObject.Section("DEFAULT")
		return section.Key(k1).String()
	}
	section := ib.iniObject.Section(k1)
	return section.Key(keys[0]).String()
}
