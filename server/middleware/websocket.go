package middleware

import (
	lg "Wave/utiles/logger"
	"net/http"
	"strings"
)

func WebSocketHeadersCheck(curlog *lg.Logger) Middleware {

	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			isContains := func(key, value string) bool {
				s1 := strings.ToLower(r.Header.Get(key))
				s2 := strings.ToLower(value)
				return s1 == s2
			}

			if isContains("Connection", "upgrade") && isContains("Upgrade", "websocket") && isContains("Sec-WebSocket-Version", "13") {

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
