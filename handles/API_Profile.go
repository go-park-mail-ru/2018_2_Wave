package handles

import (
	"Wave/server"
	"Wave/types"
	"mime"
	"mime/multipart"
	"strconv"

	"github.com/valyala/fasthttp"
)

// OnProfileGET - public API
func OnProfileGET(ctx *fasthttp.RequestCtx, sv *server.Server) {
	cookie := getSessionCookie(ctx)
	if !sv.DB.IsLoggedIn(cookie) {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		return
	}

	if profile, ok := sv.DB.GetProfile(cookie); ok {
		ctx.Write(types.Shield(profile.MarshalJSON()))
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusForbidden)
}

// OnProfilePOST - public API
func OnProfilePOST(ctx *fasthttp.RequestCtx, sv *server.Server) {
	var (
		cookie = getSessionCookie(ctx)
		user   = types.APIEditProfile{}
		form   *multipart.Form
		err    error
	)
	if !sv.DB.IsLoggedIn(cookie) {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		return
	}

	if form, err = ctx.MultipartForm(); err != nil {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		return
	}

	if err := parsForm(form, &user); err != nil {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		return
	}

	sv.DB.UpdateProfile(cookie, user)
	ctx.SetStatusCode(fasthttp.StatusCreated)
}

// OnAvatarGET - internal API
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
