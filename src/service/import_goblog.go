package service

import (
	"encoding/json"
	"github.com/Unknwon/cae/zip"
	"github.com/fuxiaohei/pugo/src/core"
	"github.com/fuxiaohei/pugo/src/model"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type (
	user struct {
		Id                        int64
		Name, Password            string
		Nick, Email, Avatar       string
		Url, Bio                  string
		CreateTime, LastLoginTime int64
		Role                      string
	}
	content struct {
		Id                int64
		Title, Slug, Text string
		Tags              []string
		CreateTime        int64
		EditTime          int64
		UpdateTime        int64
		IsComment         bool
		IsLinked          bool
		AuthorId          int64
		Template          string
		Type              string
		Status            string
		Format            string
		Comments          []*comment
		Hits              int64
	}
	comment struct {
		Id                 int64
		Author, Email, Url string
		Avatar             string
		Content            string
		CreateTime         int64
		// Content id
		Cid, Pid      int64
		Status        string
		Ip, UserAgent string
		// Is comment of admin
		IsAdmin bool
	}
	mFile struct {
		Id          int64
		Name        string
		UploadTime  int64
		Url         string
		ContentType string
		Author      int64
		IsUsed      bool
		Size        int64
		Type        string
		Hits        int64
	}
)

func (is *ImportService) importGoBlog(u *model.User, filePath string) error {
	dirPath := strings.Replace(filePath, ".zip", "_zip", -1)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}
	zip.Verbose = false
	if err := zip.ExtractTo(filePath, dirPath); err != nil {
		return err
	}

	// save contents
	if err := importGoBlogContents(u, dirPath); err != nil {
		return err
	}
	if err := importGoBlogMedia(u, dirPath); err != nil {
		return err
	}
	return nil
}

func importGoBlogContents(u *model.User, filePath string) error {
	return filepath.Walk(path.Join(filePath, "data", "content"), func(p string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(p, ".json") {
			return nil
		}
		bytes, _ := ioutil.ReadFile(p)
		content := new(content)
		if err := json.Unmarshal(bytes, content); err != nil {
			return err
		}
		// save article
		if content.Type == "article" {
			return migrateGoBlogArticle(content, u)
		}
		// save page
		if content.Type == "page" {
			return migrateGoBlogPage(content, u)
		}
		return nil
	})
}

func migrateGoBlogArticle(content *content, u *model.User) error {
	// create article object
	article := &model.Article{
		UserId:        u.Id,
		Title:         content.Title,
		Link:          content.Slug,
		Body:          content.Text,
		BodyType:      model.ARTICLE_BODY_MARKDOWN,
		TagString:     strings.Join(content.Tags, ","),
		Hits:          content.Hits,
		Status:        model.ARTICLE_STATUS_PUBLISH,
		CommentStatus: model.ARTICLE_COMMENT_OPEN,
	}
	if content.Status == "draft" {
		article.Status = model.ARTICLE_STATUS_DRAFT
	}
	if strings.Contains(article.Body, "<!--more-->") {
		article.Preview = strings.Split(article.Body, "<!--more-->")[0]
	}

	if _, err := Article.Write(article); err != nil {
		return err
	}

	// refresh time
	if _, err := core.Db.Exec("UPDATE article SET create_time = ? , update_time = ? WHERE id = ?", content.CreateTime, content.UpdateTime, article.Id); err != nil {
		return err
	}

	// save comments
	commentIds := make(map[int64]int64)
	for _, cmt := range content.Comments {
		if strings.ToLower(cmt.Status) != "approved" {
			continue
		}
		c := &model.Comment{
			Name:      cmt.Author,
			UserId:    0,
			Email:     cmt.Email,
			Url:       cmt.Url,
			AvatarUrl: cmt.Avatar,
			Body:      cmt.Content,
			From:      model.COMMENT_FROM_ARTICLE,
			FromId:    article.Id,
			ParentId:  commentIds[cmt.Pid],
			Status:    model.COMMENT_STATUS_WAIT,
			UserIp:    cmt.Ip,
			UserAgent: cmt.UserAgent,
		}
		if cmt.Email == u.Email {
			c.UserId = u.Id
		}
		c.Status = model.COMMENT_STATUS_APPROVED

		if _, err := core.Db.Insert(c); err != nil {
			return err
		}
		if _, err := core.Db.Exec("UPDATE comment SET create_time = ? WHERE id = ?", cmt.CreateTime, c.Id); err != nil {
			return err
		}

		commentIds[cmt.Id] = c.Id
	}

	// refresh comment count
	Comment.updateCommentCount(model.COMMENT_FROM_ARTICLE, article.Id)

	return nil
}

func migrateGoBlogPage(content *content, u *model.User) error {
	// create page object
	page := &model.Page{
		UserId:        u.Id,
		Title:         content.Title,
		Link:          content.Slug,
		Body:          content.Text,
		BodyType:      model.PAGE_BODY_MARKDOWN,
		Hits:          content.Hits,
		Status:        model.ARTICLE_STATUS_PUBLISH,
		CommentStatus: model.ARTICLE_COMMENT_OPEN,
		Template:      "page.tmpl",
	}

	// save page
	_, err := Page.Write(page)
	return err
}

func importGoBlogMedia(u *model.User, filePath string) error {
	file := path.Join(filePath, "data/files.json")
	bytes, _ := ioutil.ReadFile(file)
	files := make([]*mFile, 0)
	if err := json.Unmarshal(bytes, &files); err != nil {
		return err
	}
	for _, f := range files {

		ext := path.Ext(f.Url)
		fileType := Setting.Media.GetType(ext)
		if fileType == 0 {
			continue
		}

		m := &model.Media{
			Name:     f.Name,
			FileName: path.Base(f.Url),
			FilePath: f.Url,
			FileSize: f.Size,
			FileType: fileType,
			UserId:   u.Id,
		}
		if _, err := core.Db.Insert(m); err != nil {
			return err
		}
		if _, err := core.Db.Exec("UPDATE media SET create_time = ? WHERE id = ?", f.UploadTime, m.Id); err != nil {
			return err
		}

		mediaFile := path.Join(filePath, f.Url)
		os.MkdirAll(path.Dir(f.Url), os.ModePerm)

		readF, err := os.Open(mediaFile)
		if err != nil {
			return err
		}
		writeF, err := os.Create(f.Url)

		defer readF.Close()
		defer writeF.Close()

		if err != nil {
			return err
		}
		if _, err := io.Copy(writeF, readF); err != nil {
			return err
		}

	}
	return nil
}
