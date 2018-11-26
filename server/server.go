package server

import (
	"Wave/server/api"
	"Wave/server/database"
	"Wave/server/implemenation"
	mw "Wave/server/middleware"
	mc "Wave/server/metrics"
	"log"
	"Wave/session"

	_ "math/rand"

	"Wave/utiles/config"
	lg "Wave/utiles/logger"
	"net/http"
	"net"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	_ "google.golang.org/grpc/codes"

	_ "golang.org/x/net/context"
	"google.golang.org/grpc"
	_ "sync"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Start(path string) {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()
	curlog := lg.Construct()
	session.RegisterAuthCheckerServer(server, implementation.NewSessionManager(curlog))

	fmt.Println("starting server at :8081")
	go server.Serve(lis)

	conf := config.Configure(path)
	
	prof := mc.Construct()

	db := database.New(curlog)

	grcpConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()

	API := &api.Handler{
		DB: *db,
		LG: curlog,
		Prof: prof,
	}

	API.SessManager = session.NewAuthCheckerClient(grcpConn)

	r := mux.NewRouter()
	//r.HandleFunc("/metrics", promhttp.Handler().(http.HandlerFunc)).Methods("GET")

	r.HandleFunc("/", mw.Chain(API.SlashHandler))
	r.HandleFunc("/users", mw.Chain(API.RegisterPOSTHandler, mw.CORS(conf.CC, curlog))).Methods("POST") //!
	r.HandleFunc("/users/me", mw.Chain(API.MeGETHandler, mw.Auth(curlog), mw.CORS(conf.CC, curlog))).Methods("GET")
	r.HandleFunc("/users/me", mw.Chain(API.EditMePUTHandler, mw.Auth(curlog), mw.CORS(conf.CC, curlog))).Methods("PUT")
	r.HandleFunc("/users/{name}", mw.Chain(API.UserGETHandler, mw.CORS(conf.CC, curlog))).Methods("GET")
	r.HandleFunc("/users/leaders", mw.Chain(API.LeadersGETHandler, mw.CORS(conf.CC, curlog))).Methods("GET")
	r.HandleFunc("/session", mw.Chain(API.LoginPOSTHandler, mw.CORS(conf.CC, curlog))).Methods("POST") 
	r.HandleFunc("/session", mw.Chain(API.LogoutDELETEHandler, mw.Auth(curlog), mw.CORS(conf.CC, curlog))).Methods("DELETE") //!
	r.HandleFunc("/checkme", mw.Chain(API.IsLoggedIn, mw.CORS(conf.CC, curlog))).Methods("GET") //!

	r.HandleFunc("/users/me", mw.Chain(API.EditMeOPTHandler, mw.OptionsPreflight(conf.CC, curlog))).Methods("OPTIONS")
	r.HandleFunc("/session",  mw.Chain(API.LogoutOPTHandler, mw.OptionsPreflight(conf.CC, curlog))).Methods("OPTIONS")

	http.ListenAndServe(conf.SC.Port, handlers.RecoveryHandler()(r))
}