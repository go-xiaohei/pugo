package admin

import (
	"github.com/fuxiaohei/pugo/src/middle"
	"github.com/fuxiaohei/pugo/src/service"
)

type AdvBackupController struct {
	middle.AuthorizeRequire
	middle.AdminRender
	middle.Responsor
}

func (abc *AdvBackupController) Get() {
	abc.Title("BACKUP - PUGO")
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
