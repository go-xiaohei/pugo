package parser

import (
	"io"
)

var (
	_ Parser = (*MdParser)(nil)
	_ Parser = (*CommonParser)(nil)
)

type (
	// Parser defines the options of a parser
	Parser interface {
		Is([]byte) bool                         // check can be parsed
		Parse([]byte) ([]Block, error)          // parse bytes
		ParseReader(io.Reader) ([]Block, error) // parse io.Reader
	}
	Block interface {
		New() Block               // new block
		Type() string             // the type of block
		Is(typeBytes []byte) bool // is block type
		Write([]byte) error       // write to block
		Bytes() []byte            // read from block
	}
)
