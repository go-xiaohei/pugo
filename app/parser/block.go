package parser

import (
	"bytes"
	"html/template"
)

const (
	BLOCK_MARKDOWN = "markdown"
)

type (
	ContentBlock interface {
		Block
		Preview() []byte            // preview bytes
		Content() []byte            // content bytes
		PreviewHTML() template.HTML // preview HTML
		ContentHTML() template.HTML // content HTML
	}
	MarkdownBlock struct {
		// todo: implement ContentBlock and use in builder
		buf *bytes.Buffer
	}
)

// new MarkdownBlock
func (mb *MarkdownBlock) New() Block {
	return &MarkdownBlock{
		buf: bytes.NewBuffer(nil),
	}
}

// get MarkdownBlock type
func (mb *MarkdownBlock) Type() string {
	return BLOCK_MARKDOWN
}

// check is MarkdownBlock
func (mb *MarkdownBlock) Is(mark []byte) bool {
	return bytes.Equal(mark, []byte(BLOCK_MARKDOWN))
}

// write data to MarkdownBlock
func (mb *MarkdownBlock) Write(data []byte) error {
	_, err := mb.buf.Write(data)
	return err
}

// get bytes in MarkdownBlock
func (mb *MarkdownBlock) Bytes() []byte {
	return mb.buf.Bytes()
}
