package handles

import (
	"Wave/server"
	"Wave/types"
	"mime/multipart"

	"github.com/valyala/fasthttp"
)

// OnLogInPOST - public API
func OnLogInPOST(ctx *fasthttp.RequestCtx, sv *server.Server) {
	var (
		user = types.APIUser{}
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

	if sv.DB.IsSignedUp(user) {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusAccepted)
	var (
		cookieValue   = sv.DB.LogIn(user)
		sessionCookie = makeSessionCookie(cookieValue)
	)
	setCookie(ctx, sessionCookie)
}
