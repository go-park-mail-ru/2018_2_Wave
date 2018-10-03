package api

import (
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
