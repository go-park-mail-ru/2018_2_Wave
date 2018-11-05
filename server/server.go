package server

import (
	"Wave/server/api"
	"Wave/server/database"
	"Wave/utiles/config"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Start(path string) {
	conf := config.Configure(path)

	fmt.Println(conf.CC.Methods[3])

	/*
		headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
		originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	*/
	r := mux.NewRouter()

	db := database.New(conf.DC)

	API := &api.Handler{
		DB: *db,
	}

	r.HandleFunc("/", API.SlashHandler).Methods("GET")
	r.HandleFunc("/users", API.RegisterHandler).Methods("POST")
	r.HandleFunc("/users/me", API.GetMeHandler).Methods("GET")
	r.HandleFunc("/users/me", API.EditMeHandler).Methods("PUT")
	r.HandleFunc("/users/me", API.OptEditMeHandler).Methods("OPTIONS")
	r.HandleFunc("/users/{name}", API.GetUserHandler).Methods("GET")
	r.HandleFunc("/users/leaders", API.GetLeadersHandler).Methods("GET")
	r.HandleFunc("/session", API.LoginHandler).Methods("POST")
	r.HandleFunc("/session", API.LogoutHandler).Methods("DELETE")
	r.HandleFunc("/session", API.OptLogoutHandler).Methods("OPTIONS")

	log.Fatal(http.ListenAndServe(conf.SC.Port, r))
}
