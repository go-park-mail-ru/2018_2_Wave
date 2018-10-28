package user

import (
	"Wave/app/generated/models"
	"Wave/app/generated/restapi/operations/user"

	"Wave/utiles/walhalla"

	"github.com/go-openapi/runtime/middleware"
)

// walhalla:file { model:NewModel }

// walhalla:gen {}
func MyProfile(params user.MyProfileParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	// cookie := misc.GetSessionCookie(ctx)
	// if profile, ok := model.GetProfile(cookie); ok {
	// 	//ctx.Write(Must(profile.MarshalJSON()))
	// 	ctx.SetStatusCode(fasthttp.StatusOK)
	// 	return
	// }
	// ctx.SetStatusCode(fasthttp.StatusForbidden)
	return user.NewMyProfileOK().WithPayload(&models.UserExtended{
		Avatar:   "https://i.ytimg.com/vi/nc-zywmhB78/hqdefault.jpg",
		Score:    228,
		Username: "Hang",
	})
}

// walhalla:gen {}
func UpdateMyProfile(params user.UpdateMyProfileParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	// cookie := misc.GetSessionCookie(ctx)
	// model.UpdateProfile(cookie, user)
	// ctx.SetStatusCode(fasthttp.StatusAccepted)
	return middleware.NotImplemented("kek")
}

// walhalla:gen {}
func UsersProfile(params user.UsersProfileParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	// cookie := misc.GetSessionCookie(ctx)
	// model.UpdateProfile(cookie, user)
	// ctx.SetStatusCode(fasthttp.StatusAccepted)
	return middleware.NotImplemented("kek")
}
