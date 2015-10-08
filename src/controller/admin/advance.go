package admin

import (
	"github.com/fuxiaohei/pugo/src/middle"
	"github.com/fuxiaohei/pugo/src/model"
	"github.com/fuxiaohei/pugo/src/service"
	"github.com/lunny/tango"
	"os"
)

type AdvBackupController struct {
	tango.Ctx

	middle.AuthorizeRequire
	middle.AdminRender
	middle.Responsor
}

func (abc *AdvBackupController) Get() {
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
	aic.Render("advance_import.tmpl")
}
