package parser

import "bytes"

const (
	BLOCK_MARKDOWN = "markdown"
)

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
	return bytes.Equal(mark, []byte(BLOCK_MARKDOWN))
}

func (mb *MarkdownBlock) Write(data []byte) error {
	_, err := mb.buf.Write(data)
	return err
}

func (mb *MarkdownBlock) Bytes() []byte {
	return mb.buf.Bytes()
}
