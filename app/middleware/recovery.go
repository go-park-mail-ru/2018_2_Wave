package middleware

import (
	"Wave/utiles/walhalla"
	"github.com/go-openapi/runtime/middleware"
	"net/http"
)

func Recovery(next walhalla.MiddlewareFunction, ctx *walhalla.Context) walhalla.MiddlewareFunction {
	return func(r *http.Request) middleware.Responder {
		defer func() {
			// if err := recover(); err != nil {
			// 	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			// 	ctx.Response.SetBody([]byte{})
			// 	sv.Log.WithFields(logger.Fields{
			// 		"name": functionName,
			// 		"type": "handle",
			// 	}).Errorln(err)
			// }
		}()
		return next(r)
	}
}
