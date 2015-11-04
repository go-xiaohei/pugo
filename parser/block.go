package parser

import "bytes"

const (
	BLOCK_INI      = "ini"
	BLOCK_MARKDOWN = "markdown"
)

type IniBlock struct {
	data []byte
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

type MarkdownBlock struct {
	buf *bytes.Buffer
}

func (mb *MarkdownBlock) New() Block {
	return &MarkdownBlock{
		buf: bytes.NewBuffer(nil),
	}
}

func (mb *MarkdownBlock) Type() string {
	return BLOCK_MARKDOWN
}

func (mb *MarkdownBlock) Is(mark []byte) bool {
	return bytes.Equal(mark, []byte(BLOCK_INI))
}

func (mb *MarkdownBlock) Write(data []byte) error {
	_, err := mb.buf.Write(data)
	return err
}

func (mb *MarkdownBlock) Bytes() []byte {
	return mb.buf.Bytes()
}
