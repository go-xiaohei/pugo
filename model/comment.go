package model

import (
	"errors"
	"github.com/go-xiaohei/pugo-static/parser"
	"strings"
)

var (
	ErrCommentBlockWrong = errors.New("comment-blocks-wrong")
)

type Comment struct {
	Disqus *CommentDisqus `ini:"disqus"`
}

type CommentDisqus struct {
	Site string `ini:"site"`
}

func NewComment(blocks []parser.Block) (*Comment, error) {
	if len(blocks) != 1 {
		return nil, ErrCommentBlockWrong
	}
	block, ok := blocks[0].(parser.MetaBlock)
	if !ok {
		return nil, ErrCommentBlockWrong
	}
	c := new(Comment)
	// disqus
	disqus := new(CommentDisqus)
	if err := block.MapTo("disqus", disqus); err != nil {
		return nil, err
	}
	if disqus.Site != "" {
		c.Disqus = disqus
	}
	return c, nil
}

func (c *Comment) String() string {
	using := []string{}
	if c.Disqus != nil {
		using = append(using, "disqus")
	}
	return strings.Join(using, ",")
}

func (c *Comment) IsOK() bool {
	if c.Disqus != nil {
		return true
	}
	return false
}
