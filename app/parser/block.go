package parser

import (
	"bytes"
	"html/template"
)

const (
	// BlockMarkdown is markdown block type string
	BlockMarkdown = "markdown"
)

type (
	// ContentBlock defines a block to save content
	// not using yet
	ContentBlock interface {
		Block
		Preview() []byte            // preview bytes
		Content() []byte            // content bytes
		PreviewHTML() template.HTML // preview HTML
		ContentHTML() template.HTML // content HTML
	}
	// MarkdownBlock is ContentBlock of markdown
	MarkdownBlock struct {
		// todo: implement ContentBlock and use in builder
		buf *bytes.Buffer
	}
)

// New returns new MarkdownBlock
func (mb *MarkdownBlock) New() Block {
	return &MarkdownBlock{
		buf: bytes.NewBuffer(nil),
	}
}

// Type get MarkdownBlock type
func (mb *MarkdownBlock) Type() string {
	return BlockMarkdown
}

// Is checks is MarkdownBlock
func (mb *MarkdownBlock) Is(mark []byte) bool {
	return bytes.Equal(mark, []byte(BlockMarkdown))
}

// write data to MarkdownBlock
func (mb *MarkdownBlock) Write(data []byte) error {
	_, err := mb.buf.Write(data)
	return err
}

// Bytes get bytes in MarkdownBlock
func (mb *MarkdownBlock) Bytes() []byte {
	return mb.buf.Bytes()
}
