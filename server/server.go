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
	r.HandleFunc("/users", API.RegisterPOSTHandler).Methods("POST")
	r.HandleFunc("/users/me", API.MeGETHandler).Methods("GET")
	r.HandleFunc("/users/me", API.EditMePUTHandler).Methods("PUT")
	r.HandleFunc("/users/{name}", API.UserGETHandler).Methods("GET")
	r.HandleFunc("/users/leaders", API.LeadersGETHandler).Methods("GET")
	r.HandleFunc("/session", API.LoginPOSTHandler).Methods("POST")
	r.HandleFunc("/session", API.LogoutDELETEHandler).Methods("DELETE")

	r.HandleFunc("/users/me", mw.Chain(API.EditMeOPTHandler, mw.CORS(conf.CC), mw.Options())).Methods("OPTIONS")
	r.HandleFunc("/session",  mw.Chain(API.LogoutOPTHandler, mw.CORS(conf.CC), mw.Options())).Methods("OPTIONS")

	log.Fatal(http.ListenAndServe(conf.SC.Port, r))
}
