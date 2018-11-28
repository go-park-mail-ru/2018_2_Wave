package game

import (
	"Wave/internal/config"
	"Wave/internal/logger"
	"Wave/internal/metrics"
	"Wave/internal/services/auth/proto"
	mw "Wave/internal/middleware"

	"net/http"

	"google.golang.org/grpc"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

type Game struct {
	*Handler
}

func NewGame(curlog *logger.Logger, Prof *metrics.Profiler, conf config.Configuration) *Game {
	var (
		AuthManager auth.AuthClient
		g = &Game{ Handler: NewHandler(curlog, Prof, AuthManager)}
		r = mux.NewRouter()
	)
	{ // get auth manager
		Auth := conf.Auth
		grpcConn, err := grpc.Dial(
			Auth.Host+Auth.Port,
			grpc.WithInsecure(),
		)
		if err != nil {
			panic(err) //TODO::
		}
		AuthManager = auth.NewAuthClient(grpcConn)
	}

	r.HandleFunc("/conn/ws", mw.Chain(g.WSHandler, mw.WebSocketHeadersCheck(curlog, Prof), mw.CORS(conf.CC, curlog, Prof))).Methods("GET")

	server := http.Server{
		Addr: conf.Game.WsPort,
		Handler: handlers.RecoveryHandler()(r),
	}
	server.ListenAndServe()

	// TODO:: log
	// http.ListenAndServe(conf.Game.WsPort, handlers.RecoveryHandler()(r))
	// TODO:: log
	return g
}
