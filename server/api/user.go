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
	cookieValue, err := sv.DB.SignUp(user)
	if cookieValue == "" && err == nil {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		reason := []byte("{\"response\":\"user already exists\"}")
		ctx.Write(reason)

		return
	} else if cookieValue != "" && err == nil {
		sessionCookie := misc.MakeSessionCookie(cookieValue)
		ctx.SetStatusCode(fasthttp.StatusCreated)
		misc.SetCookie(ctx, sessionCookie)

		return
	} else if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	} 
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
	cookieValue, err := sv.DB.LogIn(user)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	} else if cookieValue == "" {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		reason := []byte("{\"response\":\"wrong password\"}")
		ctx.Write(reason)

		return
	} else if cookieValue != "" {
		sessionCookie := misc.MakeSessionCookie(cookieValue)
		ctx.SetStatusCode(fasthttp.StatusAccepted)
		misc.SetCookie(ctx, sessionCookie)

		return
	}
}

// OnLogOut - public API
// walhalla: {
// 		URI: 		/user/logout,
// 		Method: 	POST,
// 		Auth: 		true
// }
func OnLogOut(ctx *fasthttp.RequestCtx, sv *server.Server) {
	cookie := misc.GetSessionCookie(ctx)
	err := sv.DB.LogOut(cookie)
	
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}

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
	isLoggedIn, err := sv.DB.IsLoggedIn(cookie)
	
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	} else if isLoggedIn {
		profile, e := sv.DB.GetProfile(cookie)
		if e != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			
			return
		}
		ctx.Write(types.Must(profile.MarshalJSON()))
		ctx.SetStatusCode(fasthttp.StatusOK)

		return
	} else if !isLoggedIn {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)

		return
	}
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
	isLoggedIn, err := sv.DB.IsLoggedIn(cookie)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	} else if isLoggedIn {
		isUpdated, e := sv.DB.UpdateProfile(cookie, user)
		if e != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	
			return
		}
		if isUpdated {
			ctx.SetStatusCode(fasthttp.StatusAccepted)

			return
		} else if !isUpdated{
			ctx.SetStatusCode(fasthttp.StatusForbidden)
			reason := []byte("{\"response\":\"bad update\"}")
			ctx.Write(reason)

			return
		}
	} else if !isLoggedIn {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)

		return
	}
}

// OnAvatarGET - public API
// walhalla: {
// 		URI: 		/img/avatars/:uid,
// 		Method: 	GET,
// 		Auth: 		true
// }
func OnAvatarGET(ctx *fasthttp.RequestCtx, sv *server.Server) {
	cookie := misc.GetSessionCookie(ctx)
	isLoggedIn, err := sv.DB.IsLoggedIn(cookie)
	
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	} else if isLoggedIn {
		avatarSource, e := sv.DB.GetAvatar(cookie)
		if e != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			
			return
		}
		ctx.Write([]byte(avatarSource))
		ctx.SetStatusCode(fasthttp.StatusOK)

		return
	} else if !isLoggedIn {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)

		return
	}
}

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