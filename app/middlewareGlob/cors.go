package middlewareGlob

import (
	"fmt"

	"Wave/app/middlewareGlob/cors"
	"Wave/utiles/walhalla"
	"Wave/utiles"
)

func Cors(next walhalla.GlobalMiddlewareFunction, ctx *walhalla.Context) walhalla.GlobalMiddlewareFunction {
	Conf, ok := ctx.Config.(*utiles.MainConfig)
	if !ok {
		panic(fmt.Errorf("Unexpected config type"))
	}

	return cors.New(cors.Options{
		AllowedOrigins:   Conf.CORS.Origins,
		AllowedHeaders:   Conf.CORS.Headers,
		AllowedMethods:   Conf.CORS.Methods,
		AllowCredentials: Conf.CORS.Credentials,
	}).CorsMiddleware(next)
}
