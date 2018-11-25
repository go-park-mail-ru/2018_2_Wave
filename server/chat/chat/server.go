package chat

import (
	"Wave/services/chat/api"
	ps "Wave/server/database"
	mw "Wave/server/middleware"
	cf "Wave/utiles/config"
	lg "Wave/utiles/logger"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Start serving
func Start(confPath string) {
	var (
		curlog = lg.New()
		conf   = cf.New(confPath)
		db     = ps.New(curlog)
		API    = api.New(db)
		r      = mux.NewRouter()
	)
	r.HandleFunc("/conn/ws", mw.Chain(API.WSHandler, mw.WebSocketHeadersCheck(curlog), mw.CORS(conf.CC, curlog))).Methods("GET")

	http.ListenAndServe(conf.SC.Port, handlers.RecoveryHandler()(r))
}
