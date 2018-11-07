package server

import (
	"Wave/server/api"
	"Wave/server/database"
	mw "Wave/server/middleware"
	"Wave/utiles/config"
	//lg "Wave/utiles/logger"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

func Start(path string) {
	conf := config.Configure(path)
	r := mux.NewRouter()

	db := database.New(conf.DC)

	API := &api.Handler{
		DB: *db,
	}

	//Wavelog := lg.Construct()

	r.HandleFunc("/", mw.Chain(API.SlashHandler, mw.Auth(), mw.CORS(conf.CC))).Methods("GET")
	r.HandleFunc("/users", mw.Chain(API.RegisterPOSTHandler, mw.CORS(conf.CC))).Methods("POST")
	r.HandleFunc("/users/me", mw.Chain(API.MeGETHandler, mw.Auth(), mw.CORS(conf.CC))).Methods("GET")
	r.HandleFunc("/users/me", mw.Chain(API.EditMePUTHandler, mw.CORS(conf.CC), mw.Auth())).Methods("PUT")
	r.HandleFunc("/users/{name}", mw.Chain(API.UserGETHandler, mw.CORS(conf.CC))).Methods("GET")
	r.HandleFunc("/users/leaders", mw.Chain(API.LeadersGETHandler, mw.CORS(conf.CC))).Methods("GET")
	r.HandleFunc("/session", mw.Chain(API.LoginPOSTHandler, mw.CORS(conf.CC))).Methods("POST")
	r.HandleFunc("/session", mw.Chain(API.LogoutDELETEHandler, mw.CORS(conf.CC), mw.Auth())).Methods("DELETE")

	r.HandleFunc("/users/me", mw.Chain(API.EditMeOPTHandler, mw.OptionsPreflight(conf.CC))).Methods("OPTIONS")
	r.HandleFunc("/session",  mw.Chain(API.LogoutOPTHandler, mw.OptionsPreflight(conf.CC))).Methods("OPTIONS")

	r.HandleFunc("/conn/lobby", mw.Chain(API.LobbyHandler, mw.CORS(conf.CC), mw.WebSocketHeadersCheck())).Methods("GET")

	http.ListenAndServe(conf.SC.Port, handlers.RecoveryHandler()(r))
}
