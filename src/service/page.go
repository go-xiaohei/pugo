package service

import (
	"errors"
	"pugo/src/core"
	"pugo/src/model"
)

var (
	Page = new(PageService)

	ErrPageDisallowLink = errors.New("page-disallow-link")
)

type PageService struct{}

func (ps *PageService) Write(v interface{}) (*Result, error) {
	page, ok := v.(*model.Page)
	if !ok {
		return nil, ErrServiceFuncNeedType(ps.Write, page)
	}
	for _, str := range Setting.Content.PageDisallowLink {
		if str == page.Link {
			return nil, ErrPageDisallowLink
		}
	}
	if page.Id > 0 {
		if _, err := core.Db.Where("id = ?", page.Id).
			Cols("title,link,update_time,body,body_type,status,comment_status,top_link,template").
			Update(page); err != nil {
			return nil, err
		}
	} else {
		if _, err := core.Db.Insert(page); err != nil {
			return nil, err
		}
	}
	return newResult(ps.Write, page), nil
}
