package user

import (
	"Wave/app/generated/restapi/operations/user"
	// "Wave/app/generated/models"
	"Wave/utiles/walhalla"

	"github.com/go-openapi/runtime/middleware"
)

// walhalla:file { model:NewModel }

// walhalla:gen { mdw:[] }
func SignupUser(params user.SignupUserParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	//if !model.IsSignedUp(user.AsUser()) {
	//	var (
	//		cookieValue   = model.SignUp(user)
	//		sessionCookie = misc.MakeSessionCookie(cookieValue)
	//	)
	//	ctx.SetStatusCode(fasthttp.StatusCreated)
	//	misc.SetCookie(ctx, sessionCookie)
	//} else {
	//	ctx.SetStatusCode(fasthttp.StatusForbidden)
	//}
	return middleware.NotImplemented("kek")
}
