package middleware

import (
	"Wave/utiles/config"
	"Wave/utiles/cors"
	lg "Wave/utiles/logger"
	"net/http"
	"strings"
)

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
