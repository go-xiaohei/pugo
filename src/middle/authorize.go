package middle

import (
	"github.com/go-xiaohei/pugo/src/model"
	"github.com/go-xiaohei/pugo/src/service"
	"github.com/lunny/tango"
	"net/http"
	"time"
)

var (
	AuthTokenHttpHeader   = "PUGO-TOKEN"
	AuthTokenFormField    = "PUGO_TOKEN"
	AuthTokenCookieName   = AuthTokenFormField
	AuthFailUrl           = "/admin/login"
	AuthUserTemplateField = "AuthUser"

	_ IAuthorize = (*AuthorizeCheck)(nil)
	_ IAuthorize = (*AuthorizeRequire)(nil)
)

// authorize interface
type IAuthorize interface {
	ReadToken(*tango.Context) string
	WriteToken(*tango.Context, *model.UserToken)
	SetAuthUser(*model.User)
	OnAuthFail(*tango.Context) bool // bool returns the rest handlers should be continued
}

type AuthorizeCheck struct {
	AuthUser *model.User
}

// read token from request
func (_ *AuthorizeCheck) ReadToken(ctx *tango.Context) string {
	token := ctx.Header().Get(AuthTokenHttpHeader)
	if token == "" {
		token = ctx.Cookie(AuthTokenCookieName)
	}
	if token == "" {
		token = ctx.Form(AuthTokenFormField)
	}
	return token
}

// write token to response
func (_ *AuthorizeCheck) WriteToken(ctx *tango.Context, token *model.UserToken) {
	if token == nil {
		ctx.Cookies().Set(&http.Cookie{
			Name:    AuthTokenCookieName,
			Value:   "",
			Path:    "/",
			MaxAge:  -3600,
			Expires: time.Now().Add(-1 * time.Hour),
		})
		return
	}
	ctx.Cookies().Set(&http.Cookie{
		Name:     AuthTokenCookieName,
		Value:    token.Hash,
		Path:     "/",
		Expires:  time.Unix(token.ExpireTime, 0),
		MaxAge:   int(token.ExpireTime - time.Now().Unix()),
		HttpOnly: true,
	})
}

// set auth user
func (ac *AuthorizeCheck) SetAuthUser(user *model.User) {
	ac.AuthUser = user
}

// on auth fail handler,
// authorize checker only try to assign auth user,
// not to check access
func (_ *AuthorizeCheck) OnAuthFail(_ *tango.Context) bool {
	return false
}

// authorize require handler
type AuthorizeRequire struct {
	AuthorizeCheck
}

// if auth fail, clean token,
// redirect to login page
func (ar *AuthorizeRequire) OnAuthFail(ctx *tango.Context) bool {
	// clean token cookie
	ctx.Cookies().Set(&http.Cookie{
		Name:    AuthTokenCookieName,
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour),
		MaxAge:  -3600,
	})
	// todo : delete token
	ctx.Redirect(AuthFailUrl)
	return true
}

// authorize handler
func Authorizor() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		auth, ok := ctx.Action().(IAuthorize)
		if !ok {
			ctx.Next()
			return
		}

		token := auth.ReadToken(ctx)
		if token != "" {
			var (
				opt = service.UserVerifyOption{
					Hash:           token,
					Origin:         "webpage",
					Extend:         true,
					ExtendDuration: 3600 * 24,
				}
				user = new(model.User)
			)
			if err := service.Call(service.User.Verify, opt, user); err == nil {
				auth.SetAuthUser(user)
				if render, ok := ctx.Action().(ITheme); ok {
					render.Assign(AuthUserTemplateField, user)
				}
				ctx.Next()
				return
			}
		}

		if auth.OnAuthFail(ctx) {
			return
		}
		ctx.Next()
	}
}
