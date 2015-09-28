package service

import (
	"github.com/fuxiaohei/pugo/src/core"
	"github.com/fuxiaohei/pugo/src/model"
)

type MessageService struct{}

func (ms *MessageService) Save(v interface{}) (*Result, error) {
	m, ok := v.(*model.Message)
	if !ok {
		return nil, ErrServiceFuncNeedType(ms.Save, m)
	}
	if _, err := core.Db.Insert(m); err != nil {
		return nil, err
	}
	return newResult(ms.Save, m), nil
}
