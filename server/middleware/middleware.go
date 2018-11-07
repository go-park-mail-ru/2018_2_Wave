package middleware

import (
	"Wave/utiles/config"
	"Wave/utiles/models"
	"Wave/utiles/cors"
	//lg "Wave/utiles/logger"
	"log"

	//"Wave/utiles/misc"
	//"log"
	"strings"
	"fmt"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func CORS(CC config.CORSConfiguration) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
				originToSet := cors.SetOrigin(r.Header.Get("Origin"), CC.Origins)
				if originToSet == "" {
					rw.WriteHeader(http.StatusForbidden)
					log.Println("yeah")
					return
				}
				rw.Header().Set("Access-Control-Allow-Origin", originToSet)
				rw.Header().Set("Access-Control-Allow-Headers", strings.Join(CC.Headers, ", "))
				rw.Header().Set("Access-Control-Allow-Credentials", CC.Credentials)
				rw.Header().Set("Access-Control-Allow-Methods", strings.Join(CC.Methods, ", "))
				//log.Println("yeah")
				hf(rw, r)
			}
	}
}

func OptionsPreflight(CC config.CORSConfiguration) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
				originToSet := cors.SetOrigin(r.Header.Get("Origin"), CC.Origins)
				if originToSet == "" {
					rw.Header().Set("Access-Control-Allow-Origin", originToSet)
					rw.Header().Set("Access-Control-Allow-Headers", strings.Join(CC.Headers, ", "))
					rw.Header().Set("Access-Control-Allow-Credentials", CC.Credentials)
					rw.Header().Set("Access-Control-Allow-Methods", strings.Join(CC.Methods, ", "))
					rw.WriteHeader(http.StatusForbidden)

					return
				}
				rw.Header().Set("Access-Control-Allow-Origin", originToSet)
				rw.Header().Set("Access-Control-Allow-Headers", strings.Join(CC.Headers, ", "))
				rw.Header().Set("Access-Control-Allow-Credentials", CC.Credentials)
				rw.Header().Set("Access-Control-Allow-Methods", strings.Join(CC.Methods, ", "))
				rw.WriteHeader(http.StatusOK)
				//hf(rw, r)

				return
		}
	}
}

func Auth() Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("session")
			if err != nil {
				fr := models.ForbiddenRequest{
					Reason: "Not authorized.",
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

func WebSocketHeadersCheck() Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Connection") == "Upgrade" && r.Header.Get("Upgrade") == "websocket" && r.Header.Get("Sec-Websocket-Version") == "13" {
				hf(rw, r)
			}
			rw.WriteHeader(http.StatusExpectationFailed)
		}
	}
}

func Chain(hf http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		hf = m(hf)
	}
	return hf
}
