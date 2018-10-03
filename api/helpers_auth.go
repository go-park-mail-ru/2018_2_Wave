package api

import (
	"Wave/server"

	"github.com/valyala/fasthttp"
)

const sessionCookieLifeTime = 60 * 24 * 365
const sessionCookieName = "session"

func makeSessionCookie(value string) *fasthttp.Cookie {
	loginCookie := &fasthttp.Cookie{}
	loginCookie.SetMaxAge(sessionCookieLifeTime)
	loginCookie.SetKey(sessionCookieName)
	loginCookie.SetSecure(false)
	loginCookie.SetValue(value)
	return loginCookie
}

func getSessionCookie(ctx *fasthttp.RequestCtx) string {
	return string(ctx.Request.Header.Cookie(sessionCookieName))
}

func setCookie(ctx *fasthttp.RequestCtx, cookie *fasthttp.Cookie) {
	ctx.Response.Header.SetCookie(cookie)
}

func requireAuth(ctx *fasthttp.RequestCtx, sv *server.Server) bool {
	cookie := getSessionCookie(ctx)
	if !sv.DB.IsLoggedIn(cookie) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return true
	}
	return false
}

func requireUnAuth(ctx *fasthttp.RequestCtx, sv *server.Server) bool {
	cookie := getSessionCookie(ctx)
	if !sv.DB.IsLoggedIn(cookie) {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		return false
	}
	return true
}
