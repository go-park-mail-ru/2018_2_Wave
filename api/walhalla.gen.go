// The code was generated. Dont edit the one
//
package api

import (
	"Wave/server"
	"Wave/types"
	"github.com/valyala/fasthttp"
)

// NOTES:
// 1. insead of '&lt;'(less) the template generates '&it;'. I don't know why but it is.

func HandlerOnLogInPOST(ctx *fasthttp.RequestCtx, sv *server.Server) {

	println("request for OnLogInPOST")

	target := types.APIUser{}

	if form, err := ctx.MultipartForm(); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	} else if err := target.UnmarshalForm(form); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if !target.Validate() {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	OnLogInPOST(ctx, sv, target)

}

func HandlerOnProfileGET(ctx *fasthttp.RequestCtx, sv *server.Server) {

	println("request for OnProfileGET")

	cookie := getSessionCookie(ctx)
	if !sv.DB.IsLoggedIn(cookie) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	OnProfileGET(ctx, sv)

}

func HandlerOnProfilePOST(ctx *fasthttp.RequestCtx, sv *server.Server) {

	println("request for OnProfilePOST")

	cookie := getSessionCookie(ctx)
	if !sv.DB.IsLoggedIn(cookie) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	target := types.APIEditProfile{}

	if form, err := ctx.MultipartForm(); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	} else if err := target.UnmarshalForm(form); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if !target.Validate() {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	OnProfilePOST(ctx, sv, target)

}

func HandlerOnAvatarGET(ctx *fasthttp.RequestCtx, sv *server.Server) {

	println("request for OnAvatarGET")

	cookie := getSessionCookie(ctx)
	if !sv.DB.IsLoggedIn(cookie) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	OnAvatarGET(ctx, sv)

}

func HandlerOnSignUpPOST(ctx *fasthttp.RequestCtx, sv *server.Server) {

	println("request for OnSignUpPOST")

	target := types.APISignUp{}

	if form, err := ctx.MultipartForm(); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	} else if err := target.UnmarshalForm(form); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if !target.Validate() {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	OnSignUpPOST(ctx, sv, target)

}

func UseAPI(sv *server.Server) {

	sv.POST("/login", HandlerOnLogInPOST)

	sv.GET("/me", HandlerOnProfileGET)

	sv.POST("/user", HandlerOnProfilePOST)

	sv.GET("/img/avatars/:uid", HandlerOnAvatarGET)

	sv.POST("/register", HandlerOnSignUpPOST)

}
