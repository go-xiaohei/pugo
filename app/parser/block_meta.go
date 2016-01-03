package parser

import (
	"bytes"

	"gopkg.in/ini.v1"
)

const (
	// BlockIni is IniBlock's type
	BlockIni = "ini"
)

type (
	// MetaBlock defines how to read items from block,
	// use to read map items
	MetaBlock interface {
		Block                                      // need block interface
		MapTo(section string, v interface{}) error // map data to struct, such as json.Umashal
		MapHash(section string) map[string]string  // get k-v map in section
		Keys(section string) []string              // get keys in section
		Item(k1 string, keys ...string) string     // get item with keys
	}
)

// IniBlock parses block with ini content
type IniBlock struct {
	data      []byte
	iniObject *ini.File
}

// New returns new ini block
func (ib *IniBlock) New() Block {
	return new(IniBlock)
}

// Type get ini's block type,
// implement Block
func (ib *IniBlock) Type() string {
	return BlockIni
}

// Is checks is block type
func (ib *IniBlock) Is(mark []byte) bool {
	return bytes.Equal(mark, []byte(BlockIni))
}

// Write writes bytes to this block
func (ib *IniBlock) Write(data []byte) error {
	ib.data = append(ib.data, data...)
	ib.iniObject = nil // todo: need locker
	return nil
}

// Bytes reads bytes in this block
func (ib *IniBlock) Bytes() []byte {
	return bytes.TrimRight(ib.data, "\n")
}

// MapTo maps section data to struct
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

// MapHash gets section data as k-v map
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

// Keys gets keys in section
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

// Item gets item with keys
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
