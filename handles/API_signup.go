package handles

import (
	"Wave/server"
	"Wave/types"
	"mime/multipart"

	"github.com/valyala/fasthttp"
)

// OnSignUpPOST - public API
func OnSignUpPOST(ctx *fasthttp.RequestCtx, sv *server.Server) {
	var (
		user = types.APISignUp{}
		form *multipart.Form
		err  error
	)

	if form, err = ctx.MultipartForm(); err != nil {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		return
	}

	if err := parsForm(form, &user); err != nil {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		return
	}

	if sv.DB.IsSignedUp(user.AsAPIUser()) {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	var (
		cookieValue   = sv.DB.SignUp(user)
		sessionCookie = makeSessionCookie(cookieValue)
	)
	setCookie(ctx, sessionCookie)
}
