package api

import (
	"Wave/server"
	"Wave/types"
	"github.com/valyala/fasthttp"
)

// OnSignUpPOST - public API
// walhalla: {
// 		URI: 		/register,
// 		Method: 	POST,
// 		Data: 		form,
// 		Target:  	types.APISignUp,
// 		Validation: yes,
// 		Auth: 		any
// }
func OnSignUpPOST(ctx *fasthttp.RequestCtx, sv *server.Server, user types.APISignUp) {
	ctx.SetStatusCode(fasthttp.StatusCreated)
	var (
		cookieValue   = sv.DB.SignUp(user)
		sessionCookie = makeSessionCookie(cookieValue)
	)
	setCookie(ctx, sessionCookie)
}
