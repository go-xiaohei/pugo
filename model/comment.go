package model

import (
	"errors"
	"strings"

	"github.com/go-xiaohei/pugo-static/parser"
)

var (
	ErrCommentBlockWrong = errors.New("comment-blocks-wrong")
)

// Comment options
type Comment struct {
	Disqus *CommentDisqus `ini:"disqus"`
}

// Comment options of Disqus
type CommentDisqus struct {
	Site string `ini:"site"`
}

// blocks to Comment
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

// Comment pasred third-party comments system,
// return as disqus,duoshuo, or empty string
func (c *Comment) String() string {
	using := []string{}
	if c.Disqus != nil {
		using = append(using, "disqus")
	}
	return strings.Join(using, ",")
}

// IsOK means is comment enabled,
// not empty settings
func (c *Comment) IsOK() bool {
	if c.Disqus != nil {
		return true
	}
	return false
}
