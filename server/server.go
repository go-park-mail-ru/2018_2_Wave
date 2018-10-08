package server

import (
	"Wave/server/database"
	"Wave/utiles"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type Server struct {
	entery *fasthttp.Server
	router *fasthttprouter.Router
	Conf   *utiles.MainConfig
	DB     *database.DB
}

type RequestHandler func(*fasthttp.RequestCtx, *Server)

//*****************|

func New(pathToConf string) (sv *Server) {
	sv = &Server{
		entery: &fasthttp.Server{},
		router: fasthttprouter.New(),
		Conf:   &utiles.MainConfig{},
		DB:     database.New(),
	}
	sv.entery.Handler = sv.router.Handler
	sv.Conf.ReadFromFile(pathToConf)
	return sv
}

func (sv *Server) Start() error {
	port := sv.Conf.Server.Port
	return sv.entery.ListenAndServe(port)
}

//*****************| Handlers

func (sv *Server) wrapHandle(handle RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		handle(ctx, sv)
	}
}

func (sv *Server) GET(path string, handle RequestHandler) {
	sv.router.GET(path, sv.wrapHandle(handle))
}

func (sv *Server) POST(path string, handle RequestHandler) {
	sv.router.POST(path, sv.wrapHandle(handle))
}

func (sv *Server) PUT(path string, handle RequestHandler) {
	sv.router.PUT(path, sv.wrapHandle(handle))
}

func (sv *Server) HEAD(path string, handle RequestHandler) {
	sv.router.HEAD(path, sv.wrapHandle(handle))
}
