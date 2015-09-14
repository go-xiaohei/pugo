package service

import (
	"pugo/src/core"
	"pugo/src/model"
	"strings"
)

var (
	Article *ArticleService = new(ArticleService)
)

type ArticleService struct{}

func (as *ArticleService) Write(v interface{}) (*Result, error) {
	article, ok := v.(*model.Article)
	if !ok {
		return nil, ErrServiceFuncNeedType(article, as.Write)
	}
	if article.Id > 0 {
		if _, err := core.Db.Where("id = ?", article.Id).
			Cols("title,link,update_time,preview,body,topic,tag_string,status,comment_status").
			Update(article); err != nil {
			return nil, err
		}
	} else {
		if _, err := core.Db.Insert(article); err != nil {
			return nil, err
		}
	}
	if article.TagString != "" {
		if err := as.SaveTags(article.Id, article.TagString); err != nil {
			return nil, err
		}
	}
	return newResult(as.Write, article), nil
}

func (as *ArticleService) SaveTags(id int64, tagStr string) error {
	// delete old tags
	if _, err := core.Db.Where("article_id = ?", id).Delete(new(model.ArticleTag)); err != nil {
		return err
	}
	// save new tags
	tags := strings.Split(strings.Replace(tagStr, "，", ",", -1), ",")
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if t != "" {
			if _, err := core.Db.Insert(&model.ArticleTag{ArticleId: id, Tag: t}); err != nil {
				return err
			}
		}
	}
	return nil
}
