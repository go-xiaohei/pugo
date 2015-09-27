package service

import (
	"errors"
	"github.com/fuxiaohei/pugo/src/core"
	"github.com/fuxiaohei/pugo/src/model"
	"github.com/fuxiaohei/pugo/src/utils"
	"strings"
)

var (
	Page = new(PageService)

	ErrPageDisallowLink = errors.New("page-disallow-link")
	ErrPageNotFound     = errors.New("page-not-found")
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

type PageReadOption struct {
	Id        int64
	Link      string
	Status    int8
	IsHit     bool
	IsPublish bool
}

func (a PageReadOption) toWhereString() (string, []interface{}) {
	args := make([]interface{}, 0)
	strs := make([]string, 0)
	if a.Id > 0 {
		strs = append(strs, "id = ?")
		args = append(args, a.Id)
	}
	if a.Link != "" {
		strs = append(strs, "link = ?")
		args = append(args, a.Link)
	}
	if a.Status > 0 {
		strs = append(strs, "status = ?")
		args = append(args, a.Status)
	} else {
		strs = append(strs, "status != ?")
		args = append(args, model.ARTICLE_STATUS_DELETE)
	}
	return strings.Join(strs, " AND "), args
}

func (as *PageService) Read(v interface{}) (*Result, error) {
	opt, ok := v.(PageReadOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(as.Read, opt)
	}
	whereStr, whereArgs := opt.toWhereString()
	a := new(model.Page)
	if _, err := core.Db.Where(whereStr, whereArgs...).Get(a); err != nil {
		return nil, err
	}
	if a.Id == 0 {
		return nil, ErrPageNotFound
	}
	if opt.IsPublish && !a.IsPublish() {
		return nil, ErrPageNotFound
	}
	if opt.IsHit {
		if _, err := core.Db.Exec("UPDATE page SET hits = hits + 1 WHERE id = ?", a.Id); err != nil {
			return nil, err
		}
	}
	return newResult(as.Read, a), nil
}

type PageListOption struct {
	Status  int8
	Order   string
	Page    int
	Size    int
	IsCount bool
}

func preparePageListOption(opt PageListOption) PageListOption {
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

func (ps *PageService) List(v interface{}) (*Result, error) {
	opt, ok := v.(PageListOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(ps.List, opt)
	}
	opt = preparePageListOption(opt)

	sess := core.Db.NewSession().Limit(opt.Size, (opt.Page-1)*opt.Size).OrderBy(opt.Order)
	defer sess.Close()
	if opt.Status == 0 {
		sess.Where("status != ?", model.PAGE_STATUS_DELETE)
	} else {
		sess.Where("status = ?", opt.Status)
	}

	pages := make([]*model.Page, 0)
	if err := sess.Find(&pages); err != nil {
		return nil, err
	}
	res := newResult(ps.List, &pages)
	if opt.IsCount {
		// the session had been used, reset condition to count
		if opt.Status == 0 {
			sess.Where("status != ?", model.ARTICLE_STATUS_DELETE)
		} else {
			sess.Where("status = ?", opt.Status)
		}
		count, err := sess.Count(new(model.Article))
		if err != nil {
			return nil, err
		}
		res.Set(utils.CreatePager(opt.Page, opt.Size, int(count)))
	}
	return res, nil
}
