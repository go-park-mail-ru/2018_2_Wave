package main

import (
	"Wave/internal/config"
	"Wave/internal/database"
	lg "Wave/internal/logger"
	mc "Wave/internal/metrics"
	mw "Wave/internal/middleware"
	"Wave/internal/services/api"

	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	confPath = "./conf.json"

	logPath = "./logs/"
	logFile = "api.log"
)

//go:generate easyjson ./internal/config/
//go:generate easyjson ./internal/models/
//go generate protoc --go_out=plugins=grpc:. ./internal/services/auth/proto/*.proto

func main() {
	var (
		conf   = config.Configure(confPath)
		curlog = lg.Construct(logPath, logFile)
		prof   = mc.Construct()
		db     = database.New(curlog)
		API    = api.NewHandler(db, curlog, prof)
		r      = mux.NewRouter()
	)
	r.HandleFunc("/metrics", promhttp.Handler().(http.HandlerFunc)).Methods("GET")

	r.HandleFunc("/", mw.Chain(API.SlashHandler))
	r.HandleFunc("/users", mw.Chain(API.RegisterPOSTHandler, mw.CORS(conf.CC, curlog, prof))).Methods("POST")
	r.HandleFunc("/users/me", mw.Chain(API.MeGETHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("GET")
	r.HandleFunc("/users/me", mw.Chain(API.EditMePUTHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("PUT")
	r.HandleFunc("/users/me/locale", mw.Chain(API.EditPOSTLocale, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("POST")
	r.HandleFunc("/users/{name}", mw.Chain(API.UserGETHandler, mw.CORS(conf.CC, curlog, prof))).Methods("GET")
	r.HandleFunc("/session", mw.Chain(API.LoginPOSTHandler, mw.CORS(conf.CC, curlog, prof))).Methods("POST")
	r.HandleFunc("/session", mw.Chain(API.LogoutDELETEHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("DELETE")

	r.HandleFunc("/users/me", mw.Chain(API.EditMeOPTHandler, mw.OptionsPreflight(conf.CC, curlog, prof))).Methods("OPTIONS")
	r.HandleFunc("/session", mw.Chain(API.LogoutOPTHandler, mw.OptionsPreflight(conf.CC, curlog, prof))).Methods("OPTIONS")

	r.HandleFunc("/apps", mw.Chain(API.ShowAppsGETHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("GET")
	r.HandleFunc("/apps/categories", mw.Chain(API.ShowCategoriesGETHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("GET")
	r.HandleFunc("/apps/popular", mw.Chain(API.ShowAppsPopularGETHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("GET")
	//r.HandleFunc("/apps/{name}", mw.Chain(API.AppGETHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("GET")
	r.HandleFunc("/apps/{name}", mw.Chain(API.AppPersonalGETHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("GET")
	r.HandleFunc("/apps/category/{category}", mw.Chain(API.CategoryGETHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("GET")
	r.HandleFunc("/me/apps", mw.Chain(API.AddAppPOSTHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("POST")
	r.HandleFunc("/me/apps", mw.Chain(API.MeShowAppsGetHandler, mw.Auth(curlog, prof), mw.CORS(conf.CC, curlog, prof))).Methods("GET")

	curlog.Sugar.Infow("starting api server on "+conf.SC.Host+conf.SC.Port,
		"source", "main.go")

	http.ListenAndServe(conf.SC.Port, handlers.RecoveryHandler()(r))
}
