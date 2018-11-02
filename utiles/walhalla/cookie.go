package walhalla

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

type cookieResponder struct {
	middleware.Responder
	ctx *Context
}

func (cr cookieResponder) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	cr.Responder.WriteResponse(rw, producer)
	for _, c := range cr.ctx.outCookies {
		http.SetCookie(rw, c)
	}
}

func CookieMiddleware(next MiddlewareFunction, ctx *Context) MiddlewareFunction {
	return func(r *http.Request) middleware.Responder {
		return &cookieResponder{
			Responder: next(r),
			ctx: ctx,
		}
	}
}
