package main

import (
	"Wave/internal/service/app"
	"Wave/internal/database"
	"Wave/server/implemenation"
	mw "Wave/internal/middleware"
	mc "Wave/internal/metrics"
	"Wave/session"

	"Wave/internal/config"
	lg "Wave/internal/logger"
	"net/http"
	"net"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"

	"google.golang.org/grpc"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:generate easyjson ../internal/config/
//go:generate easyjson ../internal/models/
//go:generate go run .

func main() {
	path := "../configs/conf.json"
	conf := config.Configure(path)
	curlog := lg.Construct()
	prof := mc.Construct()

	lis, err := net.Listen("tcp", conf.GRPCC.Port)
	if err != nil {
		curlog.Sugar.Infow("can't listen on port",
		"source", "main.go",)
	}

	server := grpc.NewServer()
	session.RegisterAuthCheckerServer(server, implementation.NewSessionManager(curlog))

	fmt.Println("starting server at "+conf.GRPCC.Port)
	go server.Serve(lis)

	db := database.New(curlog)

	grcpConn, err := grpc.Dial(
		conf.GRPCC.Host+conf.GRPCC.Port,
		grpc.WithInsecure(),
	)
	if err != nil {
		curlog.Sugar.Infow("can't connect to grpc server",
		"source", "main.go",)
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

	r.HandleFunc("/users/me", mw.Chain(API.EditMeOPTHandler, mw.OptionsPreflight(conf.CC, curlog))).Methods("OPTIONS")
	r.HandleFunc("/session",  mw.Chain(API.LogoutOPTHandler, mw.OptionsPreflight(conf.CC, curlog))).Methods("OPTIONS")

	http.ListenAndServe(conf.SC.Port, handlers.RecoveryHandler()(r))
}
