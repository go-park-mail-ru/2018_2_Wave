package api

import (
	"Wave/server"
	"Wave/types"
	"github.com/valyala/fasthttp"
	"mime"
	"strconv"
)

// OnProfileHEAD - public API
// walhalla: {
// 		URI: 		/me,
// 		Method: 	HEAD,
// 		Auth: 		yes
// }
func OnProfileHEAD(ctx *fasthttp.RequestCtx, sv *server.Server) {
	// empty
}

// OnProfileGET - public API
// walhalla: {
// 		URI: 		/me,
// 		Method: 	GET,
// 		Auth: 		yes
// }
func OnProfileGET(ctx *fasthttp.RequestCtx, sv *server.Server) {
	cookie := getSessionCookie(ctx)
	if profile, ok := sv.DB.GetProfile(cookie); ok {
		ctx.Write(types.Must(profile.MarshalJSON()))
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusForbidden)
}

// OnProfilePOST - public API
// walhalla: {
// 		URI: 		/me,
// 		Method: 	POST,
// 		Data: 		form,
// 		Target:  	types.APIEditProfile,
// 		Validation: yes,
// 		Auth: 		yes
// }
func OnProfilePOST(ctx *fasthttp.RequestCtx, sv *server.Server, user types.APIEditProfile) {
	cookie := getSessionCookie(ctx)
	sv.DB.UpdateProfile(cookie, user)
	ctx.SetStatusCode(fasthttp.StatusAccepted)
}

// OnAvatarGET - public API
// walhalla: {
// 		URI: 		/img/avatars/:uid,
// 		Method: 	GET,
// 		Auth: 		yes
// }
func OnAvatarGET(ctx *fasthttp.RequestCtx, sv *server.Server) {
	if uid, err := strconv.Atoi(ctx.UserValue("uid").(string)); err == nil {
		if data, ok := sv.DB.GetAvatar(uid); ok {
			ctx.SetContentType(mime.TypeByExtension("png"))
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.Write(data)
			return
		}
	}
	ctx.SetStatusCode(fasthttp.StatusNoContent)
}
