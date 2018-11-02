package user

import (
	"Wave/app/misc"
	"Wave/app/generated/models"
	"Wave/app/generated/restapi/operations/user"

	"Wave/utiles/walhalla"

	"github.com/go-openapi/runtime/middleware"
)

// walhalla:gen
func MyProfile(params user.MyProfileParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	cookie := ctx.GetCookie("session")
	profile, err := model.GetProfile(cookie)
	if err != nil {
		return user.NewMyProfileInternalServerError()
	}

	return user.NewMyProfileOK().WithPayload(&profile)
}

// walhalla:gen
func UpdateMyProfile(params user.UpdateMyProfileParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	cookie := misc.GetSessionCookie(params.HTTPRequest)

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
	return user.NewUpdateMyProfileOK()
}

// walhalla:gen
func UsersProfile(params user.UsersProfileParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	return middleware.NotImplemented("ez")
}
