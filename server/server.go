package server

import (
	"Wave/server/api"
	"Wave/server/database"
	mw "Wave/server/middleware"

	"Wave/utiles/config"
	lg "Wave/utiles/logger"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

func Start(path string) {
	conf := config.Configure(path)
	curlog := lg.Construct()

	db := database.New(curlog)

	API := api.New(db)
	API.LG = curlog

	r := mux.NewRouter()

	r.HandleFunc("/chat", mw.Chain(API.ChatHandler, mw.WebSocketHeadersCheck(curlog), mw.CORS(conf.CC, curlog))).Methods("GET")

	http.ListenAndServe(conf.SC.Port, handlers.RecoveryHandler()(r))
}
