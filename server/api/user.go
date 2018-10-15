package api

import (
	"Wave/server"
	"Wave/server/misc"
	"Wave/server/types"
	"github.com/valyala/fasthttp"
	//"mime"
	//"strconv"
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
	cookieValue := sv.DB.SignUp(user)
	if cookieValue == "" {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		reason := []byte{"userAlreadyExists"}
		ctx.Write(reason)

		return
	} else {
		sessionCookie := misc.MakeSessionCookie(cookieValue)
		ctx.SetStatusCode(fasthttp.StatusCreated)
		misc.SetCookie(ctx, sessionCookie)

		return
	}

	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
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
	cookieValue := sv.DB.LogIn(user)
	if cookieValue == "" {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		reason := []byte{"WrongPassword"}
		ctx.Write(reason)

		return
	} else {
		sessionCookie := misc.MakeSessionCookie(cookieValue)
		ctx.SetStatusCode(fasthttp.StatusAccepted)
		misc.SetCookie(ctx, sessionCookie)

		return
	}

	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
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
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// OnProfileGet - public API
// walhalla: {
// 		URI: 		/user,
// 		Method: 	GET,
// 		Auth: 		true
// }
func OnProfileGet(ctx *fasthttp.RequestCtx, sv *server.Server) {
	cookie := misc.GetSessionCookie(ctx)
	if sv.DB.IsLoggedIn(cookie) {
		profile := sv.DB.GetProfile(cookie)
		ctx.Write(types.Must(profile.MarshalJSON()))
		ctx.SetStatusCode(fasthttp.StatusOK)

		return
	} else {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)

		return
	}
	
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
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
	if sv.DB.IsLoggedIn(cookie) {
		updated := sv.DB.UpdateProfile(cookie, user)
		if updated {
			ctx.SetStatusCode(fasthttp.StatusAccepted)

			return
		} else {
			ctx.SetStatusCode(fasthttp.StatusForbidden)
			reason := []byte{"bad"}
			ctx.Write(reason)

			return
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)

		return
	}
	
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
}

// OnAvatarGET - public API
// walhalla: {
// 		URI: 		/img/avatars/:uid,
// 		Method: 	GET,
// 		Auth: 		true
// }
/*
func OnAvatarGET(ctx *fasthttp.RequestCtx, sv *server.Server) {
	if uid, err := strconv.Atoi(ctx.UserValue("uid").(string)); err == nil {
		sv.DB.GetAvatar(uid);
		//ctx.SetContentType(mime.TypeByExtension("png"))
		ctx.SetStatusCode(fasthttp.StatusOK)
		//ctx.Write(data)
	}
	//ctx.SetStatusCode(fasthttp.StatusNoContent)
}
*/