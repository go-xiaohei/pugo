package admin

import (
	"fmt"
	"github.com/go-xiaohei/pugo/src/middle"
	"github.com/go-xiaohei/pugo/src/model"
	"github.com/go-xiaohei/pugo/src/service"
	"github.com/lunny/tango"
	"os"
	"path"
	"strings"
)

type AdvBackupController struct {
	tango.Ctx

	middle.AuthorizeRequire
	middle.AdminRender
	middle.Responsor
}

func (abc *AdvBackupController) Get() {
	if file := abc.Form("file"); file != "" {
		abc.ServeFile(file)
		return
	}

	files := make([]*model.BackupFile, 0)
	if err := service.Call(service.Backup.Files, nil, &files); err != nil {
		abc.RenderError(500, err)
		return
	}
	abc.Title("BACKUP - PUGO")
	abc.Assign("BackupFiles", files)
	abc.Render("advance_backup.tmpl")
}

func (abc *AdvBackupController) Backup() {
	opt := service.BackupOption{
		true, true, true, true,
	}
	fileName := ""
	if err := service.Call(service.Backup.Backup, opt, &fileName); err != nil {
		abc.JSONError(200, err)
		return
	}
	abc.JSON(map[string]interface{}{
		"file": fileName,
	})
}

func (abc *AdvBackupController) Delete() {
	if file := abc.Form("file"); file != "" {
		if err := os.RemoveAll(file); err != nil {
			abc.RenderError(500, err)
			return
		}
	}
	abc.Redirect(abc.Req().Referer())
}

type AdvImportController struct {
	tango.Ctx

	middle.AuthorizeRequire
	middle.AdminRender
	middle.Responsor
}

func (aic *AdvImportController) Get() {
	aic.Title("IMPORT - PUGO")
	if service.Import.IsImporting {
		aic.Assign("IsImporting", "true")
	}
	aic.Render("advance_import.tmpl")
}

func (aic *AdvImportController) Post() {
	opt := service.ImportOption{
		User: aic.AuthUser,
	}
	t := strings.ToLower(aic.Param("type"))
	if t == "goblog" {
		opt.Type = service.IMPORT_TYPE_GOBLOG

		_, h, err := aic.Req().FormFile("file")
		if err != nil {
			aic.JSONError(500, err)
			return
		}

		savePath := fmt.Sprintf("_temp/%s", h.Filename)
		if err := os.MkdirAll(path.Dir(savePath), os.ModePerm); err != nil {
			aic.JSONError(500, err)
			return
		}
		if err := aic.SaveToFile("file", savePath); err != nil {
			aic.JSONError(500, err)
			return
		}
		opt.TempFile = savePath
		if err := service.Call(service.Import.Import, opt); err != nil {
			aic.JSONError(500, err)
			return
		}
	}
	aic.JSON(nil)
}
