package middleware

import (
	"Wave/utiles/config"
	"Wave/utiles/models"
	"Wave/utiles/cors"
	lg "Wave/utiles/logger"
	"log"

	//"Wave/utiles/misc"
	//"log"
	"strings"
	"fmt"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc
type WaveLogger lg.Logger

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

func (wl *WaveLogger) Auth() Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("session")
			if err != nil {
				fr := models.ForbiddenRequest{
					Reason: "Not authorized, bitch!",
				}
		
				payload, _ := fr.MarshalJSON()
				rw.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(rw, string(payload))
		
				return
			}
			log.Println("Your cookie value is : " + c.Value)
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
