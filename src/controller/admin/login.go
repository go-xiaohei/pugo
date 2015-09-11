package admin

import (
	"github.com/lunny/tango"
	"pugo/src/middle"
	"pugo/src/model"
	"pugo/src/service"
)

type LoginController struct {
	tango.Ctx

	middle.AuthorizeCheck
	middle.AdminRender
	middle.Validator
	middle.Responsor
}

func (lc *LoginController) Get() {
	if lc.AuthUser != nil { // if the authorizeCheck find the auth user , no need to login
		lc.Redirect("/admin")
		return
	}
	lc.Title("Login to Pugo")
	lc.Render("login.tmpl")
}

type LoginForm struct {
	Username string `form:"user" binding:"Required"`
	Password string `form:"password" binding:"Required"`
}

func (lc *LoginController) Post() {
	form := new(LoginForm)
	if err := lc.Validator.Validate(form, lc); err != nil {
		lc.JSONError(200, err)
		return
	}
	var (
		opt = service.UserAuthOption{
			Name:           form.Username,
			Password:       form.Password,
			ExpireDuration: 3600 * 24 * 3,
			Origin:         "webpage",
		}
		token = new(model.UserToken)
	)
	if err := service.Call(service.User.Authorize, opt, token); err != nil {
		lc.JSONError(200, err)
		return
	}
	lc.WriteToken(lc.Context, token)
	lc.JSON(nil)
}
