package middleware

import (
	"Wave/utiles/walhalla"
	"github.com/go-openapi/runtime/middleware"
	"net/http"
)

func Cors(next walhalla.MiddlewareFunction, ctx *walhalla.Context) walhalla.MiddlewareFunction {
	// handler := cors.New(cors.Options{
	// 	AllowedOrigins:   config.Origins,
	// 	AllowedHeaders:   config.Headers,
	// 	AllowedMethods:   config.Methods,
	// 	AllowCredentials: config.Credentials,
	// })
	return func(r *http.Request) middleware.Responder {
		// handler.CorsMiddleware(func(ctx *fasthttp.RequestCtx) {
		// 	next(ctx, sv)
		// })
		return next(r)
	}
}
