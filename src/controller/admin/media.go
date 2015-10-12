package admin

import (
	"github.com/go-xiaohei/pugo/src/middle"
	"github.com/go-xiaohei/pugo/src/model"
	"github.com/go-xiaohei/pugo/src/service"
	"github.com/go-xiaohei/pugo/src/utils"
	"github.com/lunny/tango"
)

type MediaController struct {
	tango.Ctx
	middle.AuthorizeRequire
	middle.AdminRender
	middle.Responsor
}

func (mc *MediaController) Get() {
	var (
		opt        = service.MediaListOption{IsCount: true}
		mediaFiles = make([]*model.Media, 0)
		pager      = new(utils.Pager)
	)
	if err := service.Call(service.Media.List, opt, &mediaFiles, pager); err != nil {
		mc.RenderError(500, err)
		return
	}
	mc.Assign("MediaFiles", mediaFiles)
	mc.Assign("Pager", pager)
	mc.Title("MEDIA - PUGO")
	mc.Assign("MaxSize", service.Setting.Media.MaxFileSize/1024)
	mc.Render("manage_media.tmpl")
}

func (mc *MediaController) Upload() {
	opt := service.MediaUploadOption{
		Ctx:      mc.Ctx,
		User:     mc.AuthUser.Id,           // media's owner int
		FormName: mc.Form("field", "file"), // form field name
	}
	if mc.Form("from") == "editor" {
		mc.uploadFormEditor(opt)
		return
	}
	if err := service.Call(service.Media.Upload, opt); err != nil {
		mc.JSONError(500, err)
		return
	}
	mc.JSON(nil)
}

func (mc *MediaController) uploadFormEditor(opt service.MediaUploadOption) {
	m := new(model.Media)
	if err := service.Call(service.Media.Upload, opt, m); err != nil {
		mc.JSONRaw(map[string]interface{}{
			"success": 0,
			"message": err.Error(),
		})
		return
	}
	mc.JSONRaw(map[string]interface{}{
		"success": 1,
		"url":     "/" + m.FilePath,
	})
}

type MediaDeleteController struct {
	tango.Ctx
	middle.AuthorizeRequire
	middle.AdminRender
}

func (mdc *MediaDeleteController) Get() {
	id := mdc.FormInt64("id")
	if id > 0 {
		if err := service.Call(service.Media.Delete, id); err != nil {
			mdc.RenderError(500, err)
			return
		}
	}
	mdc.Redirect(mdc.Req().Referer())
}
