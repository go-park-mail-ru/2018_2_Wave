package main

import (
	"Wave/internal/config"
	"Wave/internal/database"
	lg "Wave/internal/logger"
	mc "Wave/internal/metrics"
	mw "Wave/internal/middleware"
	gm "Wave/internal/services/game"

	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	confPath = "./conf.json"

	logPath = "./logs/"
	logFile = "game.log"
)

func main() {
	var (
		conf   = config.Configure(confPath)
		log = lg.Construct(logPath, logFile)
		db     = database.New(log)
		prof   = mc.Construct()
		g      = gm.NewHandler(log, prof, db)
		r      = mux.NewRouter()
	)
	r.HandleFunc("/conn/ws", mw.Chain(g.WSHandler, mw.WebSocketHeadersCheck(log, prof), mw.CORS(conf.CC, log, prof))).Methods("GET")
	http.ListenAndServe(conf.Game.WsPort, handlers.RecoveryHandler()(r))
}
