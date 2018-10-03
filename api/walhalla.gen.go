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

func HandlerOnLeaderbordGET(ctx *fasthttp.RequestCtx, sv *server.Server) {

	println("request for OnLeaderbordGET")

	ctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	cookie := getSessionCookie(ctx)
	if !sv.DB.IsLoggedIn(cookie) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	OnLeaderbordGET(ctx, sv)

}

func HandlerOnLogInPOST(ctx *fasthttp.RequestCtx, sv *server.Server) {

	println("request for OnLogInPOST")

	ctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

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

func HandlerOnLogOutGET(ctx *fasthttp.RequestCtx, sv *server.Server) {

	println("request for OnLogOutGET")

	ctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	OnLogOutGET(ctx, sv)

}

func HandlerOnProfileHEAD(ctx *fasthttp.RequestCtx, sv *server.Server) {

	println("request for OnProfileHEAD")

	ctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	cookie := getSessionCookie(ctx)
	if !sv.DB.IsLoggedIn(cookie) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	OnProfileHEAD(ctx, sv)

}

func HandlerOnProfileGET(ctx *fasthttp.RequestCtx, sv *server.Server) {

	println("request for OnProfileGET")

	ctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	cookie := getSessionCookie(ctx)
	if !sv.DB.IsLoggedIn(cookie) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	OnProfileGET(ctx, sv)

}

func HandlerOnProfilePOST(ctx *fasthttp.RequestCtx, sv *server.Server) {

	println("request for OnProfilePOST")

	ctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

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

	ctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	cookie := getSessionCookie(ctx)
	if !sv.DB.IsLoggedIn(cookie) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	OnAvatarGET(ctx, sv)

}

func HandlerOnSignUpPOST(ctx *fasthttp.RequestCtx, sv *server.Server) {

	println("request for OnSignUpPOST")

	ctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

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

	sv.GET("/users/:start/:count", HandlerOnLeaderbordGET)

	sv.POST("/login", HandlerOnLogInPOST)

	sv.GET("/logout", HandlerOnLogOutGET)

	sv.HEAD("/me", HandlerOnProfileHEAD)

	sv.GET("/me", HandlerOnProfileGET)

	sv.POST("/me", HandlerOnProfilePOST)

	sv.GET("/img/avatars/:uid", HandlerOnAvatarGET)

	sv.POST("/signup", HandlerOnSignUpPOST)

}
