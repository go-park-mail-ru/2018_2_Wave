package middleware

import (
	lg "Wave/utiles/logger"
	"net/http"
)

func WebSocketHeadersCheck(curlog *lg.Logger) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Connection") == "Upgrade" &&
				r.Header.Get("Upgrade") == "websocket" &&
				r.Header.Get("Sec-Websocket-Version") == "13" {

				curlog.Sugar.Infow(
					"websocket headers check succeded",
					"source", "middleware.go",
					"who", "WebSocketHeadersCheck")

				hf(rw, r)
			}
			rw.WriteHeader(http.StatusExpectationFailed)

			curlog.Sugar.Infow("websocket headers check failed",
				"source", "middleware.go",
				"who", "WebSocketHeadersCheck")

			return
		}
	}
}
