package service

import (
	"errors"
	"fmt"
	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo/src/core"
	"github.com/go-xiaohei/pugo/src/model"
	"github.com/go-xiaohei/pugo/src/utils"
	"strings"
	"time"
)

var (
	Article *ArticleService = new(ArticleService)

	ErrArticleNotFound = errors.New("article-not-found")
)

type ArticleService struct{}

func (as *ArticleService) Write(v interface{}) (*Result, error) {
	article, ok := v.(*model.Article)
	if !ok {
		return nil, ErrServiceFuncNeedType(as.Write, article)
	}

	var isUpdate = article.Id > 0
	if isUpdate {
		if _, err := core.Db.Where("id = ?", article.Id).
			Cols("title,link,update_time,preview,body,body_type,topic,tag_string,status,comment_status").
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

	defer as.msgWrite(isUpdate, article)

	return newResult(as.Write, article), nil
}

func (as *ArticleService) msgWrite(isUpdate bool, article *model.Article) {
	data := map[string]string{
		"type":   fmt.Sprint(model.MESSAGE_TYPE_ARTICLE_CREATE),
		"author": article.User().Name,
		"link":   article.Href(),
		"title":  article.Title,
		"time":   utils.TimeUnixFormat(article.CreateTime, "01/02 15:04:05"),
	}
	var body string
	if isUpdate {
		data["type"] = fmt.Sprint(model.MESSAGE_TYPE_ARTICLE_UPDATE)
		data["time"] = utils.TimeUnixFormat(article.UpdateTime, "01/02 15:04:05")
		body = com.Expand(MessageArticleUpdateTemplate, data)
	} else {
		body = com.Expand(MessageArticleCreateTemplate, data)
	}
	message := &model.Message{
		UserId:     article.UserId,
		From:       model.MESSAGE_FROM_ARTICLE,
		FromId:     article.Id,
		Type:       model.MESSAGE_TYPE_ARTICLE_CREATE,
		Body:       body,
		CreateTime: article.CreateTime,
	}
	if isUpdate {
		message.Type = model.MESSAGE_TYPE_ARTICLE_UPDATE
	}
	Message.Save(message)
}

func (as *ArticleService) Delete(v interface{}) (*Result, error) {
	id, ok := v.(int64)
	if !ok {
		return nil, ErrServiceFuncNeedType(as.Delete, id)
	}

	if _, err := core.Db.Exec("UPDATE article SET status = ? WHERE id = ?", model.ARTICLE_STATUS_DELETE, id); err != nil {
		return nil, err
	}

	defer as.msgDelete(id)

	return nil, nil
}

func (as *ArticleService) msgDelete(id int64) {
	article := new(model.Article)
	if _, err := core.Db.Where("id = ?", id).Get(article); err != nil {
		return
	}
	if article == nil || article.Id != id {
		return
	}
	data := map[string]string{
		"type":   fmt.Sprint(model.MESSAGE_TYPE_ARTICLE_REMOVE),
		"author": article.User().Name,
		"title":  article.Title,
		"time":   utils.TimeUnixFormat(article.CreateTime, "01/02 15:04:05"),
	}
	body := com.Expand(MessageArticleRemoveTemplate, data)
	message := &model.Message{
		UserId:     article.UserId,
		From:       model.MESSAGE_FROM_ARTICLE,
		FromId:     article.Id,
		Type:       model.MESSAGE_TYPE_ARTICLE_REMOVE,
		Body:       body,
		CreateTime: article.CreateTime,
	}
	Message.Save(message)
}

func (as *ArticleService) SaveTags(id int64, tagStr string) error {
	if err := as.RemoveTags(id); err != nil {
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

func (as *ArticleService) RemoveTags(id int64) error {
	// delete old tags
	if _, err := core.Db.Where("article_id = ?", id).Delete(new(model.ArticleTag)); err != nil {
		return err
	}
	return nil
}

type ArticleListOption struct {
	Status   int8
	Order    string
	Page     int
	Size     int
	IsCount  bool
	ReadTime int64
}

func prepareArticleListOption(opt ArticleListOption) ArticleListOption {
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

func (as *ArticleService) List(v interface{}) (*Result, error) {
	opt, ok := v.(ArticleListOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(as.List, opt)
	}
	opt = prepareArticleListOption(opt)

	sess := core.Db.NewSession().Limit(opt.Size, (opt.Page-1)*opt.Size).OrderBy(opt.Order)
	defer sess.Close()
	if opt.Status == 0 {
		sess.Where("status != ?", model.ARTICLE_STATUS_DELETE)
	} else {
		sess.Where("status = ?", opt.Status)
	}

	articles := make([]*model.Article, 0)
	if err := sess.Find(&articles); err != nil {
		return nil, err
	}
	if opt.ReadTime > 0 {
		for _, a := range articles {
			a.IsNewRead = (a.UpdateTime - opt.ReadTime) >= -3600
		}
	} else {
		// set first one as new article
		if len(articles) > 0 {
			articles[0].IsNewRead = true
		}
	}
	res := newResult(as.List, &articles)
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

type ArticleReadOption struct {
	Id        int64
	Link      string
	Status    int8
	IsHit     bool
	IsPublish bool
}

func (a ArticleReadOption) toWhereString() (string, []interface{}) {
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

func (as *ArticleService) Read(v interface{}) (*Result, error) {
	opt, ok := v.(ArticleReadOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(as.Read, opt)
	}
	whereStr, whereArgs := opt.toWhereString()
	a := new(model.Article)
	if _, err := core.Db.Where(whereStr, whereArgs...).Get(a); err != nil {
		return nil, err
	}
	if a.Id == 0 {
		return nil, ErrArticleNotFound
	}
	if opt.IsPublish && !a.IsPublish() {
		return nil, ErrArticleNotFound
	}
	if opt.IsHit {
		if _, err := core.Db.Exec("UPDATE article SET hits = hits + 1 WHERE id = ?", a.Id); err != nil {
			return nil, err
		}
	}
	return newResult(as.Read, a), nil
}

func (as *ArticleService) ToPublish(v interface{}) (*Result, error) {
	idPtr, ok := v.(*int64)
	if !ok {
		return nil, ErrServiceFuncNeedType(as.ToPublish, idPtr)
	}
	if _, err := core.Db.Exec("UPDATE article SET status = ? WHERE id = ?", model.ARTICLE_STATUS_PUBLISH, *idPtr); err != nil {
		return nil, err
	}
	return nil, nil
}

type ArchiveListOption struct {
	Order string
	Limit int
}

func prepareArchiveListOption(opt ArchiveListOption) ArchiveListOption {
	if opt.Order == "" {
		opt.Order = "create_time DESC"
	}
	if opt.Limit < 1 {
		opt.Limit = 100
	}
	return opt
}

func (as *ArticleService) Archive(v interface{}) (*Result, error) {
	opt, ok := v.(ArchiveListOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(as.Archive, opt)
	}
	opt = prepareArchiveListOption(opt)
	archives := make([]*model.ArticleArchive, 0)
	if err := core.Db.Where("status = ?", model.ARTICLE_STATUS_PUBLISH).OrderBy("create_time DESC").Find(&archives); err != nil {
		return nil, err
	}
	oldYear := 0
	for _, a := range archives {
		t := time.Unix(a.CreateTime, 0)
		if oldYear != t.Year() {
			oldYear = t.Year()
			a.IsNewYear = true
		}
	}
	return newResult(as.Archive, &archives), nil
}
