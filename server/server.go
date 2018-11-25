package server

import (
	"Wave/server/api"
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
	r.HandleFunc("/", mw.Chain(API.SlashHandler)).Methods("GET")
	r.HandleFunc("/users", mw.Chain(API.RegisterPOSTHandler, mw.CORS(conf.CC, curlog))).Methods("POST")
	r.HandleFunc("/users/me", mw.Chain(API.MeGETHandler, mw.Auth(curlog), mw.CORS(conf.CC, curlog))).Methods("GET")
	r.HandleFunc("/users/me", mw.Chain(API.EditMePUTHandler, mw.Auth(curlog), mw.CORS(conf.CC, curlog))).Methods("PUT")
	r.HandleFunc("/users/{name}", mw.Chain(API.UserGETHandler, mw.CORS(conf.CC, curlog))).Methods("GET")
	r.HandleFunc("/users/leaders", mw.Chain(API.LeadersGETHandler, mw.CORS(conf.CC, curlog))).Methods("GET")
	r.HandleFunc("/session", mw.Chain(API.LoginPOSTHandler, mw.CORS(conf.CC, curlog))).Methods("POST")
	r.HandleFunc("/session", mw.Chain(API.LogoutDELETEHandler, mw.Auth(curlog), mw.CORS(conf.CC, curlog))).Methods("DELETE")

	r.HandleFunc("/users/me", mw.Chain(API.EditMeOPTHandler, mw.OptionsPreflight(conf.CC, curlog))).Methods("OPTIONS")
	r.HandleFunc("/session", mw.Chain(API.LogoutOPTHandler, mw.OptionsPreflight(conf.CC, curlog))).Methods("OPTIONS")

	r.HandleFunc("/conn/ws", mw.Chain(API.WSHandler, mw.WebSocketHeadersCheck(curlog), mw.CORS(conf.CC, curlog))).Methods("GET")

	println("started")
	http.ListenAndServe(conf.SC.Port, handlers.RecoveryHandler()(r))
}
