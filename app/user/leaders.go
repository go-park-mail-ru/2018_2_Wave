package user

import (
	"Wave/app/generated/restapi/operations/user"
	// "Wave/app/generated/models"
	"Wave/utiles/walhalla"

	"github.com/go-openapi/runtime/middleware"
)

// walhalla:file { model:NewModel }

// walhalla:gen {}
func Leaders(params user.LeadersParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	// cookie := misc.GetSessionCookie(ctx)
	// model.UpdateProfile(cookie, user)
	// ctx.SetStatusCode(fasthttp.StatusAccepted)
	return middleware.NotImplemented("kek")
}
