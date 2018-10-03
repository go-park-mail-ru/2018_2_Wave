package api

import (
	"Wave/server"
	"Wave/types"
	"github.com/valyala/fasthttp"
)

// OnLogInPOST - public API
// walhalla: {
// 		URI: 		/login,
// 		Method: 	POST,
// 		Data: 		form,
// 		Target:  	types.APIUser,
// 		Validation: yes,
// 		Auth: 		any
// }
func OnLogInPOST(ctx *fasthttp.RequestCtx, sv *server.Server, user types.APIUser) {
	ctx.SetStatusCode(fasthttp.StatusAccepted)
	var (
		cookieValue   = sv.DB.LogIn(user)
		sessionCookie = makeSessionCookie(cookieValue)
	)
	setCookie(ctx, sessionCookie)
}

// OnLogOutGET - public API
// walhalla: {
// 		URI: 		/logout,
// 		Method: 	GET,
// 		Auth: 		any
// }
func OnLogOutGET(ctx *fasthttp.RequestCtx, sv *server.Server) {
	cookie := getSessionCookie(ctx)
	sv.DB.LogOut(cookie)
}
