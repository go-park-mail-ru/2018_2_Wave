package server

import (
	"Wave/server/api"
	"Wave/server/database"
	"Wave/utiles/config"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Init(path string) http.Server {
	conf := config.Configure(path)

	r := SetMuxAPI(conf)
	server := http.Server{
		Addr:    conf.SC.Port,
		Handler: r,
	}

	return server
}

func Start(server http.Server) error {
	log.Println("starting server at", server.Addr)
	return server.ListenAndServe()
}

func SetMuxAPI(conf config.Configuration) *mux.Router {

	r := mux.NewRouter()

	db := database.New(conf.DC)

	API := &api.Handler{
		DB: *db,
	}

	r.HandleFunc("/users", API.RegisterHandler).Methods("POST")
	r.HandleFunc("/users/me", API.GetMeHandler).Methods("GET")
	r.HandleFunc("/users/me", API.EditMeHandler).Methods("PUT")
	r.HandleFunc("/users/{name}", API.GetUserHandler).Methods("GET")
	r.HandleFunc("/users/leaders", API.GetLeadersHandler).Methods("GET")
	r.HandleFunc("/session", API.LoginHandler).Methods("POST")
	r.HandleFunc("/session", API.LogoutHandler).Methods("DELETE")

	return r
}
