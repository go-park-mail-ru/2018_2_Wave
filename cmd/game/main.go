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
		curlog = lg.Construct(logPath, logFile)
		db     = database.New(curlog)
		prof   = mc.Construct()
		g      = gm.NewHandler(curlog, prof, db)
		r      = mux.NewRouter()
	)
	r.HandleFunc("/conn/ws", mw.Chain(g.WSHandler, mw.WebSocketHeadersCheck(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("GET")
	http.ListenAndServe(conf.Game.WsPort, handlers.RecoveryHandler()(r))
}
