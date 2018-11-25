package server

import (
	"Wave/server/api"
	"Wave/server/database"
	mw "Wave/server/middleware"
	mc "Wave/server/metrics"

	"Wave/utiles/config"
	lg "Wave/utiles/logger"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Start(path string) {
	conf := config.Configure(path)
	curlog := lg.Construct()
	prof := mc.Construct()

	db := database.New(curlog)

	API := &api.Handler{
		DB: *db,
		LG: curlog,
		Prof: prof,
	}

	r := mux.NewRouter()
	r.HandleFunc("/metrics", promhttp.Handler().(http.HandlerFunc)).Methods("GET")

	r.HandleFunc("/", mw.Chain(API.SlashHandler))
	r.HandleFunc("/users", mw.Chain(API.RegisterPOSTHandler, mw.CORS(conf.CC, curlog))).Methods("POST")
	r.HandleFunc("/users/me", mw.Chain(API.MeGETHandler, mw.Auth(curlog), mw.CORS(conf.CC, curlog))).Methods("GET")
	r.HandleFunc("/users/me", mw.Chain(API.EditMePUTHandler, mw.Auth(curlog), mw.CORS(conf.CC, curlog))).Methods("PUT")
	r.HandleFunc("/users/{name}", mw.Chain(API.UserGETHandler, mw.CORS(conf.CC, curlog))).Methods("GET")
	r.HandleFunc("/users/leaders", mw.Chain(API.LeadersGETHandler, mw.CORS(conf.CC, curlog))).Methods("GET")
	r.HandleFunc("/session", mw.Chain(API.LoginPOSTHandler, mw.CORS(conf.CC, curlog))).Methods("POST")
	r.HandleFunc("/session", mw.Chain(API.LogoutDELETEHandler, mw.Auth(curlog), mw.CORS(conf.CC, curlog))).Methods("DELETE")

	r.HandleFunc("/users/me", mw.Chain(API.EditMeOPTHandler, mw.OptionsPreflight(conf.CC, curlog))).Methods("OPTIONS")
	r.HandleFunc("/session",  mw.Chain(API.LogoutOPTHandler, mw.OptionsPreflight(conf.CC, curlog))).Methods("OPTIONS")

	r.HandleFunc("/conn/lobby", mw.Chain(API.LobbyHandler, mw.WebSocketHeadersCheck(curlog), mw.CORS(conf.CC, curlog))).Methods("GET")

	r.HandleFunc("/conn/lobby", mw.Chain(API.LobbyHandler, mw.WebSocketHeadersCheck(curlog), mw.CORS(conf.CC, curlog))).Methods("GET")

	http.ListenAndServe(conf.SC.Port, handlers.RecoveryHandler()(r))
}