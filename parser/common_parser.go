package parser

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

const (
	COMMON_PARSER_PREFIX = "-----"
)

type (
	// CommonParser is default parser for single content
	CommonParser struct {
		blocks []Block
	}
)

// NewCommonParser retuns new CommonParser
func NewCommonParser() *CommonParser {
	return &CommonParser{
		// set IniBlock and MarkdownBlock as default blocks to parse
		blocks: []Block{new(IniBlock), new(MarkdownBlock)},
	}
}

// check bytes can be parsed by CommonParser
func (cp *CommonParser) Is(data []byte) bool {
	data = bytes.TrimLeft(data, "\n")
	return bytes.HasPrefix(data, []byte(COMMON_PARSER_PREFIX))
}

// detect block to save parsed bytes
func (cp *CommonParser) Detect(mark []byte) Block {
	for _, b := range cp.blocks {
		if b.Is(mark) {
			return b.New()
		}
	}
	return nil
}

// parse bytes to blocks
func (cp *CommonParser) Parse(src []byte) ([]Block, error) {
	if src == nil || len(src) == 0 {
		return nil, nil
	}
	buf := bytes.NewBuffer(src)
	return cp.ParseReader(buf)
}

// parser bytes in reader to blocks
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
			if currentBlock = cp.Detect(bytes.TrimLeft(lineData, COMMON_PARSER_PREFIX)); currentBlock == nil {
				return nil, errors.New("block-parse-first-error")
			}
			continue
		}

		if bytes.HasPrefix(lineData, []byte(COMMON_PARSER_PREFIX)) {
			// try to switch
			newBlock := cp.Detect(bytes.TrimLeft(lineData, COMMON_PARSER_PREFIX))
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
	// do not forget last block
	if currentBlock != nil {
		blocks = append(blocks, currentBlock)
	}
	return blocks, nil
}
