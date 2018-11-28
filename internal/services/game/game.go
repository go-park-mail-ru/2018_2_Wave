package game

import (
	"Wave/internal/logger"
	"Wave/internal/metrics"
	"Wave/internal/services/auth/proto"
	mw "Wave/internal/middleware"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

type Game struct {
	*Handler
}

func NewGame(curlog *logger.Logger, Prof *metrics.Profiler, AuthManager auth.AuthClient) *Game {
	var (
		g = &Game{ Handler: NewHandler(curlog, Prof, AuthManager)}
		r = mux.NewRouter()
	)
	r.HandleFunc("/conn/ws", mw.Chain(g.WSHandler, mw.WebSocketHeadersCheck(curlog))).Methods("GET") // TODO:: cors
	// TODO:: log
	http.ListenAndServe(":9605", handlers.RecoveryHandler()(r))
	// TODO:: log
	return g
}
