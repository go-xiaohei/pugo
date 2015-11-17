package parser

import (
	"io"
)

var (
	_ Parser = (*MdParser)(nil)
	_ Parser = (*CommonParser)(nil)
)

type (
	Parser interface {
		Is([]byte) bool
		Parse([]byte) ([]Block, error)
		ParseReader(io.Reader) ([]Block, error)
	}
	Block interface {
		New() Block
		Type() string
		Is(typeBytes []byte) bool
		Write([]byte) error
		Bytes() []byte
	}
)
