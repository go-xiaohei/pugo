package service

import (
	"errors"
	"github.com/go-xiaohei/pugo/src/model"
)

const (
	IMPORT_TYPE_GOBLOG = iota + 1
)

var (
	Import = new(ImportService)

	ErrImportIsWorking = errors.New("import-is-working")
)

type ImportService struct {
	IsImporting bool
}

type ImportOption struct {
	TempFile string
	Type     int
	User     *model.User
}

func (is *ImportService) Import(v interface{}) (*Result, error) {
	opt, ok := v.(ImportOption)
	if !ok {
		return nil, ErrServiceFuncNeedType(is.Import, opt)
	}

	if is.IsImporting {
		return nil, ErrImportIsWorking
	}

	is.IsImporting = true
	defer func() {
		is.IsImporting = false
	}()

	if opt.Type == IMPORT_TYPE_GOBLOG {
		if err := is.importGoBlog(opt.User, opt.TempFile); err != nil {
			return nil, err
		}
		return nil, nil
	}
	return nil, nil
}
