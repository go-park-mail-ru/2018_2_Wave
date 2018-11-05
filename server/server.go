package server

import (
	"Wave/server/api"
	"Wave/server/database"
	mw "Wave/server/middleware"
	"Wave/utiles/config"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Start(path string) {
	conf := config.Configure(path)
	r := mux.NewRouter()

	db := database.New(conf.DC)

	API := &api.Handler{
		DB: *db,
	}

	r.HandleFunc("/", mw.Chain(API.SlashHandler, mw.CORS(conf.CC))).Methods("GET")
	r.HandleFunc("/users", mw.Chain(API.RegisterPOSTHandler, mw.CORS(conf.CC))).Methods("POST")
	r.HandleFunc("/users/me", mw.Chain(API.MeGETHandler, mw.CORS(conf.CC))).Methods("GET")
	r.HandleFunc("/users/me", mw.Chain(API.EditMePUTHandler, mw.CORS(conf.CC))).Methods("PUT")
	r.HandleFunc("/users/{name}", mw.Chain(API.UserGETHandler, mw.CORS(conf.CC))).Methods("GET")
	r.HandleFunc("/users/leaders", mw.Chain(API.LeadersGETHandler, mw.CORS(conf.CC))).Methods("GET")
	r.HandleFunc("/session", mw.Chain(API.LoginPOSTHandler, mw.CORS(conf.CC))).Methods("POST")
	r.HandleFunc("/session", mw.Chain(API.LogoutDELETEHandler, mw.CORS(conf.CC))).Methods("DELETE")

	r.HandleFunc("/users/me", mw.Chain(API.EditMeOPTHandler, mw.CORS(conf.CC), mw.Options())).Methods("OPTIONS")
	r.HandleFunc("/session",  mw.Chain(API.LogoutOPTHandler, mw.CORS(conf.CC), mw.Options())).Methods("OPTIONS")

	log.Fatal(http.ListenAndServe(conf.SC.Port, r))
}
