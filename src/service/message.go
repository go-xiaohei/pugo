package service

import (
	"github.com/fuxiaohei/pugo/src/core"
	"github.com/fuxiaohei/pugo/src/model"
)

var (
	Message = new(MessageService)

	MessageArticleCreateTemplate = ` <div class="message msg-{type}">
                <i class="fa fa-file-text"></i>
                <span class="author"><strong>{author}</strong></span>
                write new article
                <a href="{link}" class="article"><i>{title}
                    <span class="page-link">({link})</span>
                </i>
                </a>
                <span class="time">{time}</span>
            </div>`
	MessageArticleUpdateTemplate = ` <div class="message msg-{type}">
                <i class="fa fa-file-text"></i>
                <span class="author"><strong>{author}</strong></span>
                update article
                <a href="{link}" class="article"><i>{title}
                    <span class="page-link">({link})</span>
                </i>
                </a>
                <span class="time">{time}</span>
            </div>`
	MessageArticleRemoveTemplate = `<div class="message msg-{type}">
                <i class="fa fa-file-text"></i>
                <span class="author"><strong>{author}</strong></span>
                remove article
                <a href="#"><i>{title}</i></a>
                <span class="time">{time}</span>
            </div>`
	MessagePageCreateTemplate = ` <div class="message msg-{type}">
                <i class="fa fa-file-text"></i>
                <span class="author"><strong>{author}</strong></span>
                create new page
                <a href="{link}" class="article"><i>{title}
                    <span class="page-link">({link})</span>
                </i>
                </a>
                <span class="time">{time}</span>
            </div>`
	MessagePageUpdateTemplate = ` <div class="message msg-{type}">
                <i class="fa fa-file-text"></i>
                <span class="author"><strong>{author}</strong></span>
                update page
                <a href="{link}" class="article"><i>{title}
                    <span class="page-link">({link})</span>
                </i>
                </a>
                <span class="time">{time}</span>
            </div>`
	MessagePageRemoveTemplate = `<div class="message msg-{type}">
                <i class="fa fa-file-text"></i>
                <span class="author"><strong>{author}</strong></span>
                remove page
                <a href="#"><i>{title}</i></a>
                <span class="time">{time}</span>
            </div>`
	MessageCommentLeaveTemplate = ` <div class="message msg-{type}">
                <i class="fa fa-comments"></i>
                <span class="time">{time}</span>
                <span class="author"><strong>{author}</strong>
                    <span class="email-link">({site})</span>
                </span>
                leaves a comment to
                <a><i>{title}</i></a>
                <div class="content">{body}</div>
            </div>`
	MessageCommentReplyTemplate = ` <div class="message msg-{type}">
                <i class="fa fa-comments"></i>
                <span class="time">{time}</span>
                <span class="author"><strong>{author}</strong>
                    <span class="email-link">({site})</span>
                </span>
                leaves a comment to
                <a><i>{title}</i></a>
                <div class="content">{body}</div>
                <div class="quote">
                    <span class="user">@ {parent}</span>
                    <span class="content">{parent_content}</span>
                </div>
            </div>`
	MessageCommentRemoveTemplate = ` <div class="message msg-{type}">
                <i class="fa fa-comments"></i>
                <span class="time">{time}</span>
                <span class="author"><strong>{author}'s</strong>
                    <span class="email-link">({site})</span>
                </span>
                comment is removed in
                <a><i>{title}</i></a>
            </div>`
)

type MessageService struct{}

func (ms *MessageService) Save(v interface{}) (*Result, error) {
	m, ok := v.(*model.Message)
	if !ok {
		return nil, ErrServiceFuncNeedType(ms.Save, m)
	}
	if _, err := core.Db.Insert(m); err != nil {
		return nil, err
	}
	return newResult(ms.Save, m), nil
}
