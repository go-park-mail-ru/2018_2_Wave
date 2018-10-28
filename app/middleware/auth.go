package middleware

import (
	// "Wave/app/misc"
	"Wave/utiles/walhalla"
	"github.com/go-openapi/runtime/middleware"
	"net/http"
)

func AuthTrue(next walhalla.MiddlewareFunction, ctx *walhalla.Context) walhalla.MiddlewareFunction {
	return func(r *http.Request) middleware.Responder {
		// cookie := misc.GetSessionCookie(r)
		// if !sv.DB.IsLoggedIn(cookie) {
		// 	ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		// 	return
		// }
		return next(r)
	}
}
