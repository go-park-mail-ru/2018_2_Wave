package middleware

import (
	lg "Wave/utiles/logger"
	"Wave/utiles/models"
	"fmt"
	"net/http"
)

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
