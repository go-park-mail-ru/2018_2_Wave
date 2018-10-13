package api

import (
	"Wave/server"
	"Wave/server/misc"
	"Wave/server/types"
	"github.com/valyala/fasthttp"
	//"mime"
	"strconv"
)


// OnSignUp - public API
// walhalla: {
// 		URI: 		/user/signup,
// 		Method: 	POST,
// 		Data: 		form,
// 		Target:  	types.SignUp,
// 		Validation: true
// }
func OnSignUp(ctx *fasthttp.RequestCtx, sv *server.Server, user types.SignUp) {
	ctx.SetStatusCode(fasthttp.StatusCreated)
	var (
		cookieValue   = sv.DB.SignUp(user)
		sessionCookie = misc.MakeSessionCookie(cookieValue)
	)
	misc.SetCookie(ctx, sessionCookie)
}

// OnLogIn - public API
// walhalla: {
// 		URI: 		/user/login,
// 		Method: 	POST,
// 		Data: 		form,
// 		Target:  	types.User,
// 		Validation: true
// }
func OnLogIn(ctx *fasthttp.RequestCtx, sv *server.Server, user types.User) {
	ctx.SetStatusCode(fasthttp.StatusAccepted)
	var (
		cookieValue   = sv.DB.LogIn(user)
		sessionCookie = misc.MakeSessionCookie(cookieValue)
	)
	misc.SetCookie(ctx, sessionCookie)
}

// OnLogOut - public API
// walhalla: {
// 		URI: 		/user/logout,
// 		Method: 	POST,
// 		Auth: 		true
// }
func OnLogOut(ctx *fasthttp.RequestCtx, sv *server.Server) {
	cookie := misc.GetSessionCookie(ctx)
	sv.DB.LogOut(cookie)
}

// OnProfileGet - public API
// walhalla: {
// 		URI: 		/user,
// 		Method: 	GET,
// 		Auth: 		true
// }
func OnProfileGet(ctx *fasthttp.RequestCtx, sv *server.Server) {
	cookie := misc.GetSessionCookie(ctx)
	profile := sv.DB.GetProfile(cookie)
	ctx.Write(types.Must(profile.MarshalJSON()))
	ctx.SetStatusCode(fasthttp.StatusOK)
	//ctx.SetStatusCode(fasthttp.StatusForbidden)
}

// OnProfileEdit - public API
// walhalla: {
// 		URI: 		/user/edit,
// 		Method: 	POST,
// 		Data: 		form,
// 		Target:  	types.EditProfile,
// 		Validation: true,
// 		Auth: 		true
// }
func OnProfileEdit(ctx *fasthttp.RequestCtx, sv *server.Server, user types.EditProfile) {
	cookie := misc.GetSessionCookie(ctx)
	sv.DB.UpdateProfile(cookie, user)
	ctx.SetStatusCode(fasthttp.StatusAccepted)
}

// OnAvatarGET - public API
// walhalla: {
// 		URI: 		/img/avatars/:uid,
// 		Method: 	GET,
// 		Auth: 		true
// }
func OnAvatarGET(ctx *fasthttp.RequestCtx, sv *server.Server) {
	if uid, err := strconv.Atoi(ctx.UserValue("uid").(string)); err == nil {
		sv.DB.GetAvatar(uid);
		//ctx.SetContentType(mime.TypeByExtension("png"))
		ctx.SetStatusCode(fasthttp.StatusOK)
		//ctx.Write(data)
	}
	//ctx.SetStatusCode(fasthttp.StatusNoContent)
}
