package admin

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/xsrf"
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
	lc.Title("LOGIN - PUGO")
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

func (lc *LoginController) Logout() {
	if token := lc.ReadToken(lc.Context); token != "" {
		if err := service.Call(service.User.Unauthorize, &token); err != nil {
			lc.RenderError(500, err)
			return
		}
	}
	lc.WriteToken(lc.Context, nil)
	lc.Redirect("/admin/login")
}

type ProfileController struct {
	xsrf.Checker

	middle.AuthorizeRequire
	middle.AdminRender
	middle.Validator
	middle.Responsor
}

func (pc *ProfileController) Get() {
	pc.Title("PROFILE - PUGO")
	pc.Assign("XsrfHTML", pc.XsrfFormHtml())
	pc.Render("profile.tmpl")
}

type ProfileForm struct {
	Username string `form:"name" binding:"Required"`
	UserNick string `form:"nick" binding:"Required"`
	Email    string `form:"email" binding:"Required;Email"`
	Url      string `form:"url" binding:"Url"`
	Profile  string `form:"bio"`
}

func (f ProfileForm) toUser() *model.User {
	return &model.User{
		Name:    f.Username,
		Email:   f.Email,
		Nick:    f.UserNick,
		Profile: f.Profile,
		Url:     f.Url,
	}
}

func (pc *ProfileController) Post() {
	form := new(ProfileForm)
	if err := pc.Validator.Validate(form, pc); err != nil {
		pc.JSONError(200, err)
		return
	}
	user := form.toUser()
	user.Id = pc.AuthUser.Id
	if err := service.Call(service.User.Profile, user); err != nil {
		pc.JSONError(200, err)
		return
	}
	pc.JSON(nil)
}

type PasswordForm struct {
	Old     string `form:"old" binding:"Required"`
	New     string `form:"new" binding:"Required;MinSize(6)"`
	Confirm string `form:"confirm" binding:"Required;MinSize(6)"`
}

func (f PasswordForm) toOption() service.UserPasswordOption {
	return service.UserPasswordOption{
		Old:     f.Old,
		New:     f.New,
		Confirm: f.Confirm,
	}
}

func (pc *ProfileController) Password() {
	form := new(PasswordForm)
	if err := pc.Validator.Validate(form, pc); err != nil {
		pc.JSONError(200, err)
		return
	}
	opt := form.toOption()
	opt.User = pc.AuthUser
	if err := service.Call(service.User.Password, opt); err != nil {
		pc.JSONError(200, err)
		return
	}
	pc.JSON(nil)
}
