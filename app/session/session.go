package session

import (
	"Wave/app/generated/restapi/operations/session"
	"Wave/app/generated/models"
	"Wave/app/misc"
	"Wave/utiles/walhalla"

	"github.com/go-openapi/runtime/middleware"
)

// walhalla:file { model:NewModel }

// walhalla:gen { auth:false }
func LoginUser(params session.LoginUserParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	cookie, err := model.LogIn(*params.Body) //!
	
	if err != nil {
		return session.NewLoginUserInternalServerError()
	} else if cookie == "" {
		return session.NewLoginUserUnauthorized().WithPayload(&models.ForbiddenRequest{
			Reason: "Incorrect password.",
		})
	}

	sessionCookie := misc.MakeSessionCookie(cookie)
	misc.SetCookie(ctx, sessionCookie) //!

	return session.NewLoginUserOK()
}

// walhalla:gen { auth:true }
func LogoutUser(params session.LogoutUserParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	cookie := misc.GetSessionCookie(ctx)
	err := model.LogOut(cookie) //!
	
	if err != nil {
		return session.NewLogoutUserInternalServerError()
	}

	return session.NewLogoutUserOK()
}
