package middlewareGlob

import (
	"net/http"
	"time"

	"Wave/utiles/walhalla"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.status = status
}

func Logger(next walhalla.GlobalMiddlewareFunction, ctx *walhalla.Context) walhalla.GlobalMiddlewareFunction {
	return func(rw http.ResponseWriter, r *http.Request) {
		srw := statusResponseWriter{
			ResponseWriter: rw,
		}

		defer func(start time.Time) {
			duration := time.Since(start).Nanoseconds() / int64(time.Microsecond)
			ctx.Log.WithFields(walhalla.Fields{
				"url":      r.URL.Path,
				"type":     "handle",
				"duration": duration,
			}).Info(srw.status)
		}(time.Now())

		next(&srw, r)
	}
}
