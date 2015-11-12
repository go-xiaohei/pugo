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
		Meta(section string, v interface{}) error
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
	return nil
}

func (ib *IniBlock) Bytes() []byte {
	return bytes.TrimRight(ib.data, "\n")
}

func (ib *IniBlock) Meta(section string, v interface{}) error {
	if ib.iniObject == nil {
		var err error
		ib.iniObject, err = ini.Load(ib.Bytes())
		if err != nil {
			return err
		}
	}
	return ib.iniObject.Section(section).MapTo(v)
}
