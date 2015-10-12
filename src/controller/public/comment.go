package public

import (
	"github.com/fuxiaohei/pugo/src/middle"
	"github.com/fuxiaohei/pugo/src/model"
	"github.com/fuxiaohei/pugo/src/service"
	"github.com/lunny/tango"
	"github.com/tango-contrib/session"
	"github.com/tango-contrib/xsrf"
)

type CommentController struct {
	tango.Ctx
	xsrf.Checker
	session.Session

	middle.Validator
	middle.Responsor
}

type CommentForm struct {
	Name     string `form:"name" binding:"Required"`
	Email    string `form:"email" binding:"Required;Email"`
	Url      string `form:"url" binding:"Url"`
	Content  string `form:"content" binding:"Required"`
	ParentId int64  `form:"parent"`
	UserId   int64  `form:"uid"`
	Type     string `form:"-"`
	Id       int64  `form:"-"`
}

func (cf CommentForm) toCreateOption() service.CommentCreateOption {
	return service.CommentCreateOption{
		cf.Name, cf.Email, cf.Url,
		cf.Content, cf.ParentId, cf.UserId,
		cf.Type, cf.Id,
	}
}

func (cc *CommentController) Post() {
	form := new(CommentForm)
	if err := cc.Validator.Validate(form, cc); err != nil {
		cc.JSONError(200, err)
		return
	}
	form.Type = cc.Param("type")
	form.Id = cc.ParamInt64("id")

	// create comment object
	var (
		c = new(model.Comment)
	)
	if err := service.Call(service.Comment.Create, form.toCreateOption(), c); err != nil {
		cc.JSONError(200, err)
		return
	}

	// save comment object
	if err := service.Call(service.Comment.Save, c, c); err != nil {
		cc.JSONError(200, err)
		return
	}

	cc.JSON(map[string]interface{}{
		"comment": model.NewFrontComment(c),
	})
}
