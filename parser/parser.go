package parser

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

const (
	BLOCK_PREFIX = "-----"
)

type (
	Parser interface {
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
	CommonParser struct {
		blocks []Block
	}
)

func NewCommonParser() *CommonParser {
	return &CommonParser{
		blocks: []Block{new(IniBlock), new(MarkdownBlock)},
	}
}

func (cp *CommonParser) Detect(mark []byte) Block {
	for _, b := range cp.blocks {
		if b.Is(mark) {
			return b
		}
	}
	return nil
}

func (cp *CommonParser) ParseReader(r io.Reader) ([]Block, error) {
	var (
		currentBlock Block   = nil
		blocks       []Block = nil
		reader               = bufio.NewReader(r)
	)
	for {
		lineData, _, err := reader.ReadLine()
		// first block
		if currentBlock == nil {
			if len(lineData) == 0 {
				continue
			}
			if currentBlock = cp.Detect(bytes.TrimLeft(lineData, BLOCK_PREFIX)); currentBlock == nil {
				return nil, errors.New("block-parse-first-error")
			}
			continue
		}

		if bytes.HasPrefix(lineData, BLOCK_PREFIX) {
			// try to switch
			newBlock := cp.Detect(bytes.TrimLeft(lineData, BLOCK_PREFIX))
			if newBlock != nil {
				blocks = append(blocks, currentBlock)
				currentBlock = newBlock
				continue
			}
		}

		// write block
		if err := currentBlock.Write(append(lineData, []byte("\n")...)); err != nil {
			return nil, err
		}

		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
	}
	return blocks, nil
}
