package middleware

import (
	"Wave/internal/config"
	"Wave/internal/cors"
	lg "Wave/internal/logger"
	"Wave/internal/models"

	"fmt"
	"net/http"
	"strings"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func CORS(CC config.CORSConfiguration, curlog *lg.Logger) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			originToSet := cors.SetOrigin(r.Header.Get("Origin"), CC.Origins)
			if originToSet == "" {
				rw.WriteHeader(http.StatusForbidden)

				curlog.Sugar.Infow(
					"CORS failed",
					"source", "middleware.go",
					"who", "CORS",
				)

				return
			}
			rw.Header().Set("Access-Control-Allow-Origin", originToSet)
			rw.Header().Set("Access-Control-Allow-Headers", strings.Join(CC.Headers, ", "))
			rw.Header().Set("Access-Control-Allow-Credentials", CC.Credentials)
			rw.Header().Set("Access-Control-Allow-Methods", strings.Join(CC.Methods, ", "))

			curlog.Sugar.Infow(
				"CORS succeded",
				"source", "middleware.go",
				"who", "CORS",
			)

			hf(rw, r)
		}
	}
}

func OptionsPreflight(CC config.CORSConfiguration, curlog *lg.Logger) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			originToSet := cors.SetOrigin(r.Header.Get("Origin"), CC.Origins)
			if originToSet == "" {
				rw.Header().Set("Access-Control-Allow-Origin", originToSet)
				rw.Header().Set("Access-Control-Allow-Headers", strings.Join(CC.Headers, ", "))
				rw.Header().Set("Access-Control-Allow-Credentials", CC.Credentials)
				rw.Header().Set("Access-Control-Allow-Methods", strings.Join(CC.Methods, ", "))
				rw.WriteHeader(http.StatusForbidden)

				curlog.Sugar.Infow(
					"preflight failed",
					"source", "middleware.go",
					"who", "OptionsPreflight",
				)

				return
			}

			rw.Header().Set("Access-Control-Allow-Origin", originToSet)
			rw.Header().Set("Access-Control-Allow-Headers", strings.Join(CC.Headers, ", "))
			rw.Header().Set("Access-Control-Allow-Credentials", CC.Credentials)
			rw.Header().Set("Access-Control-Allow-Methods", strings.Join(CC.Methods, ", "))
			rw.WriteHeader(http.StatusOK)

			curlog.Sugar.Infow(
				"preflight succeded",
				"source", "middleware.go",
				"who", "OptionsPreflight",
			)

			return
		}
	}
}

func Auth(curlog *lg.Logger) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")

			if err != nil || cookie.Value == "" {
				fr := models.ForbiddenRequest{
					Reason: "Not authorized.",
				}

				payload, _ := fr.MarshalJSON()
				rw.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(rw, string(payload))

				curlog.Sugar.Infow(
					"auth check failed",
					"source", "middleware.go",
					"who", "Auth",
				)

				return
			}

			curlog.Sugar.Infow(
				"auth check succeded",
				"source", "middleware.go",
				"who", "Auth",
			)

			hf(rw, r)
		}
	}
}

func WebSocketHeadersCheck(curlog *lg.Logger) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Connection") == "Upgrade" &&
			r.Header.Get("Upgrade") == "websocket" &&
			r.Header.Get("Sec-Websocket-Version") == "13" {

				curlog.Sugar.Infow("websocket headers check succeded",
					"source", "middleware.go",
					"who", "WebSocketHeadersCheck")

				hf(rw, r)

				return
			}
			rw.WriteHeader(http.StatusExpectationFailed)

			curlog.Sugar.Infow("websocket headers check failed",
				"source", "middleware.go",
				"who", "WebSocketHeadersCheck")

			return
		}
	}
}

func Chain(hf http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		hf = m(hf)
	}
	return hf
}
