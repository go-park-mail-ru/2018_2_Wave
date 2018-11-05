package middleware

import (
	"Wave/utiles/config"
	"Wave/utiles/cors"
	"net/http"
	"strings"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func CORS(CC config.CORSConfiguration) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			originToSet := cors.SetOrigin(r.Header.Get("Origin"), CC.Origins)
			if originToSet == "" {
				rw.WriteHeader(http.StatusForbidden)
				hf(rw, r)
			}
			rw.Header().Set("Access-Control-Allow-Origin", originToSet)
			rw.Header().Set("Access-Control-Allow-Headers", strings.Join(CC.Headers, ", "))
			rw.Header().Set("Access-Control-Allow-Credentials", CC.Credentials)
			rw.Header().Set("Access-Control-Allow-Methods", strings.Join(CC.Methods, ", "))
			hf(rw, r)
		}
	}
}

func Options() Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
			hf(rw, r)
		}
	}
}

func Recovery() Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			
			hf(rw, r)
		}
	}
}


func Chain(hf http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		hf = m(hf)
	}
	return hf
}
