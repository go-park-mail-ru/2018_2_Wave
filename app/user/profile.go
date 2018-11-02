package user

import (
	"Wave/app/generated/models"
	"Wave/app/generated/restapi/operations/user"

	"Wave/utiles/walhalla"

	"github.com/go-openapi/runtime/middleware"
)

// walhalla:file { model:NewModel }

// walhalla:gen { auth:true }
func MyProfile(params user.MyProfileParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	cookie := misc.GetSessionCookie(ctx) //!
	profile, err := model.GetProfile(cookie)
	if err != nil {
		return user.NewMyProfileInternalServerError()
	}

	return user.NewMyProfileOK().WithPayload(&profile)
}

// walhalla:gen { auth:true }
func UpdateMyProfile(params user.UpdateMyProfileParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	cookie := misc.GetSessionCookie(ctx) //!

	isUpdated, err := model.UpdateProfile(*params.Body, cookie)
		if err != nil {
			return user.NewUpdateMyProfileInternalServerError()
		}
		if isUpdated {
			return user.NewUpdateMyProfileOK()
		}
		if !isUpdated {
			return user.NewUpdateMyProfileForbidden().WithPayload(&models.ForbiddenRequest{
				Reason: "Bad update parameters.",
			})
		}
}

// walhalla:gen {}
func UsersProfile(params user.UsersProfileParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	
}
