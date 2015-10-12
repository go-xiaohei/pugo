package service

import (
	"errors"
	"fmt"
	"github.com/Unknwon/com"
	"github.com/fuxiaohei/pugo/src/core"
	"github.com/fuxiaohei/pugo/src/model"
	"github.com/fuxiaohei/pugo/src/utils"
	"github.com/lunny/tango"
	"mime/multipart"
	"os"
	"path"
	"time"
)

var (
	Media = new(MediaService)

	ErrMediaTooLarge     = errors.New("media-too-large")
	ErrMediaDisAllowType = errors.New("media-disallowed-type")
)

type MediaService struct{}

// media upload option
type MediaUploadOption struct {
	Ctx      tango.Ctx
	User     int64  // media's owner int
	FormName string // form field name
}

func (ms *MediaService) Upload(v interface{}) (*Result, error) {
	opt, ok := v.(MediaUploadOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(ms.Upload, opt)
	}

	f, h, err := opt.Ctx.Req().FormFile(opt.FormName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// check file size
	size, err := getUploadFileSize(f)
	if err != nil {
		return nil, err
	}
	if (size / 1024) > Setting.Media.MaxFileSize {
		return nil, ErrMediaTooLarge
	}

	// check file type
	ext := path.Ext(h.Filename)
	fileType := Setting.Media.GetType(ext)
	if fileType == 0 {
		return nil, ErrMediaDisAllowType
	}

	// hash file name, make dir
	now := time.Now()
	hashName := utils.Md5String(fmt.Sprintf("%d%s%d", opt.User, h.Filename, now.UnixNano())) + ext
	fileName := path.Join("static/upload", hashName)
	fileDir := path.Dir(fileName)

	if !com.IsDir(fileDir) {
		if err = os.MkdirAll(fileDir, os.ModePerm); err != nil {
			return nil, err
		}
	}
	if err = opt.Ctx.SaveToFile(opt.FormName, fileName); err != nil {
		return nil, err
	}

	// save media data
	media := &model.Media{
		UserId:   opt.User,
		Name:     h.Filename,
		FileName: hashName,
		FilePath: fileName,
		FileSize: size,
		FileType: fileType,
	}
	if _, err := core.Db.Insert(media); err != nil {
		return nil, err
	}

	defer ms.msgUpload(media)

	return newResult(ms.Upload, media), nil
}

func (ms *MediaService) msgUpload(m *model.Media) {
	user, err := getUserBy("id", m.UserId)
	if err != nil {
		return
	}
	data := map[string]string{
		"type":   fmt.Sprint(model.MESSAGE_TYPE_MEDIA_UPLOAD),
		"time":   utils.TimeUnixFormat(m.CreateTime, "01/02 15:04:05"),
		"author": user.Name,
		"file":   m.Name,
	}
	message := &model.Message{
		UserId:     m.UserId,
		From:       model.MESSAGE_FROM_MEDIA,
		FromId:     m.Id,
		Type:       model.MESSAGE_TYPE_MEDIA_UPLOAD,
		CreateTime: m.CreateTime,
		Body:       com.Expand(MessageMediaUploadTemplate, data),
	}
	Message.Save(message)
}

// an interface to check Size() method
type fileSizer interface {
	Size() int64
}

// get file size
func getUploadFileSize(f multipart.File) (int64, error) {
	// if return *http.sectionReader, it is alias to *io.SectionReader
	if s, ok := f.(fileSizer); ok {
		return s.Size(), nil
	}
	// or *os.File
	if fp, ok := f.(*os.File); ok {
		fi, err := fp.Stat()
		if err != nil {
			return 0, err
		}
		return fi.Size(), nil
	}
	return 0, nil
}

type MediaListOption struct {
	Type    int
	Order   string
	Page    int
	Size    int
	IsCount bool
}

func prepareMediaListOption(opt MediaListOption) MediaListOption {
	if opt.Order == "" {
		opt.Order = "create_time DESC"
	}
	if opt.Size == 0 {
		opt.Size = 10
	}
	if opt.Page < 1 {
		opt.Page = 1
	}
	return opt
}

func (ms *MediaService) List(v interface{}) (*Result, error) {
	opt, ok := v.(MediaListOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(ms.List, opt)
	}
	opt = prepareMediaListOption(opt)

	sess := core.Db.NewSession().Limit(opt.Size, (opt.Page-1)*opt.Size).OrderBy(opt.Order)
	defer sess.Close()
	if opt.Type > 0 {
		sess.Where("file_type = ?", opt.Type)
	}

	mediaFiles := make([]*model.Media, 0)
	if err := sess.Find(&mediaFiles); err != nil {
		return nil, err
	}

	res := newResult(ms.List, &mediaFiles)

	if opt.IsCount {
		if opt.Type > 0 {
			sess.Where("file_type = ?", opt.Type)
		}
		count, err := sess.Count(new(model.Media))
		if err != nil {
			return nil, err
		}
		res.Set(utils.CreatePager(opt.Page, opt.Size, int(count)))
	}
	return res, nil
}

func (ms *MediaService) Delete(v interface{}) (*Result, error) {
	id, ok := v.(int64)
	if !ok {
		return nil, ErrServiceFuncNeedType(ms.Delete, id)
	}
	if _, err := core.Db.Exec("DELETE FROM media WHERE id = ?", id); err != nil {
		return nil, err
	}
	return nil, nil
}
