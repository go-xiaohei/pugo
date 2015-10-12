package service

import (
	"github.com/go-xiaohei/pugo/src/core"
	"github.com/go-xiaohei/pugo/src/model"
	"github.com/go-xiaohei/pugo/src/utils"
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
	MessageMediaUploadTemplate = ` <div class="message msg-{type}">
                <i class="fa fa-file-image-o"></i>
                <span class="time">{time}</span>
                <span class="author"><strong>{author}</strong></span>
                upload file
                <a href="/admin/manage/media"><i>{file}</i></a>
            </div>`
	MessageBackupCreateTemplate = `<div class="message msg-{type}">
                <i class="fa fa-file-zip-o"></i>
                all site data are backup to
                <a href="/admin/advance/backup"><i>{file}</i></a>
                <span class="time">{time}</span>
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

type MessageListOption struct {
	Page, Size int
	Order      string
	IsCount    bool
}

func prepareMessageListOption(opt MessageListOption) MessageListOption {
	if opt.Order == "" {
		opt.Order = "create_time DESC"
	}
	if opt.Page < 1 {
		opt.Page = 1
	}
	if opt.Size == 0 {
		opt.Size = 10
	}
	return opt
}

func (ms *MessageService) List(v interface{}) (*Result, error) {
	opt, ok := v.(MessageListOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(ms.List, opt)
	}
	opt = prepareMessageListOption(opt)
	msgs := make([]*model.Message, 0)
	if err := core.Db.Limit(opt.Size, (opt.Page-1)*opt.Size).OrderBy(opt.Order).Find(&msgs); err != nil {
		return nil, err
	}
	res := newResult(ms.List, &msgs)
	if opt.IsCount {
		count, err := core.Db.Count(new(model.Message))
		if err != nil {
			return nil, err
		}
		res.Set(utils.CreatePager(opt.Page, opt.Size, int(count)))
	}
	return res, nil
}
