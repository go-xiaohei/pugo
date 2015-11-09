package model

import (
	"errors"
	"github.com/go-xiaohei/pugo-static/parser"
	"gopkg.in/ini.v1"
	"strings"
)

var (
	ErrCommentBlockWrong = errors.New("comment-blocks-wrong")
)

type Comment struct {
	Disqus *CommentDisqus
}

type CommentDisqus struct {
	Site string
}

func NewComment(blocks []parser.Block) (*Comment, error) {
	if len(blocks) != 1 {
		return nil, ErrCommentBlockWrong
	}
	iniF, err := ini.Load(blocks[0].Bytes())
	if err != nil {
		return nil, err
	}
	c := new(Comment)
	// disqus
	section := iniF.Section("disqus")
	if site := section.Key("site").String(); site != "" {
		disqus := &CommentDisqus{
			Site: site,
		}
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
