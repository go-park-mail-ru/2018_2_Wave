package main

import (
	"Wave/internal/config"
	"Wave/internal/database"
	lg "Wave/internal/logger"
	mc "Wave/internal/metrics"
	mw "Wave/internal/middleware"
	"Wave/internal/services/api"
	auth "Wave/internal/services/auth/proto"

	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const (
	confPath = "./conf.json"

	logPath = "./build/logs/"
	logFile = "api.log"
)

//go:generate easyjson ./internal/config/
//go:generate easyjson ./internal/models/
//go generate protoc --go_out=plugins=grpc:. ./internal/services/auth/proto/*.proto

func main() {
	conf := config.Configure(confPath)
	curlog := lg.Construct(logPath, logFile)
	prof := mc.Construct()
	db := database.New(curlog)

	grpcConn, err := grpc.Dial(
		conf.AC.Host+conf.AC.Port,
		grpc.WithInsecure(),
	)

	if err != nil {

		curlog.Sugar.Infow("can't connect to grpc server",
			"source", "main.go")

	}

	defer grpcConn.Close()

	API := &api.Handler{
		DB:          *db,
		LG:          curlog,
		Prof:        prof,
		AuthManager: auth.NewAuthClient(grpcConn),
	}

	r := mux.NewRouter()
	r.HandleFunc("/metrics", promhttp.Handler().(http.HandlerFunc)).Methods("GET")

	r.HandleFunc("/", mw.Chain(API.SlashHandler))
	r.HandleFunc("/users", mw.Chain(API.RegisterPOSTHandler, mw.CORS(conf.CC, curlog, prof))).Methods("POST")
	r.HandleFunc("/users/me", mw.Chain(API.MeGETHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("GET")
	r.HandleFunc("/users/me", mw.Chain(API.EditMePUTHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("PUT")
	r.HandleFunc("/users/{name}", mw.Chain(API.UserGETHandler, mw.CORS(conf.CC, curlog, prof))).Methods("GET")
	r.HandleFunc("/users/leaders", mw.Chain(API.LeadersGETHandler, mw.CORS(conf.CC, curlog, prof))).Methods("GET")
	r.HandleFunc("/session", mw.Chain(API.LoginPOSTHandler, mw.CORS(conf.CC, curlog, prof))).Methods("POST")
	r.HandleFunc("/session", mw.Chain(API.LogoutDELETEHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("DELETE")

	r.HandleFunc("/users/me", mw.Chain(API.EditMeOPTHandler, mw.OptionsPreflight(conf.CC, curlog, prof))).Methods("OPTIONS")
	r.HandleFunc("/session", mw.Chain(API.LogoutOPTHandler, mw.OptionsPreflight(conf.CC, curlog, prof))).Methods("OPTIONS")
	r.HandleFunc("/me/apps", mw.Chain(API.DeleteAppOPTHandler, mw.OptionsPreflight(conf.CC, curlog, prof))).Methods("OPTIONS")

	r.HandleFunc("/apps", mw.Chain(API.ShowAppsGETHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("GET")
	r.HandleFunc("/me/apps", mw.Chain(API.AddAppPOSTHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("POST")
	r.HandleFunc("/me/apps", mw.Chain(API.MeShowAppsGetHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("GET")
	r.HandleFunc("/me/apps", mw.Chain(API.DeleteAppDELETEHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("DELETE")
	r.HandleFunc("/apps/popular", mw.Chain(API.ShowAppsPopularGETHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("GET")
	r.HandleFunc("/apps/{name}", mw.Chain(API.AppGETHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("GET")

	r.HandleFunc("/me/apps/timer", mw.Chain(API.AppTimerPOSTHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("POST")

	curlog.Sugar.Infow("starting api server on "+conf.SC.Host+conf.SC.Port,
		"source", "main.go")

	http.ListenAndServe(conf.SC.Port, handlers.RecoveryHandler()(r))
}
