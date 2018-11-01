package user

import (
	"Wave/app/generated/restapi/operations/user"
	// "Wave/app/generated/models"
	"Wave/utiles/walhalla"

	"github.com/go-openapi/runtime/middleware"
)

// walhalla:file { model:NewModel }

// walhalla:gen { auth:false }
func LoginUser(params user.LoginUserParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	// if !model.IsSignedUp(user) {
	// 	var (
	// 		cookieValue   = model.LogIn(user)
	// 		sessionCookie = misc.MakeSessionCookie(cookieValue)
	// 	)
	// 	ctx.SetStatusCode(fasthttp.StatusAccepted)
	// 	misc.SetCookie(ctx, sessionCookie)
	// } else {
	// 	ctx.SetStatusCode(fasthttp.StatusForbidden)
	// }
	return middleware.NotImplemented("kek")
}

// walhalla:gen {}
func LogoutUser(params user.LogoutUserParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	return middleware.NotImplemented("kek")
}
