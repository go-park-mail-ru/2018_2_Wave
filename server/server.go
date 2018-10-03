package server

import (
	"Wave/database"
	"fmt"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type Server struct {
	entery *fasthttp.Server
	router *fasthttprouter.Router
	DB     *database.DB
}

type RequestHandler func(*fasthttp.RequestCtx, *Server)

//*****************|

func New() (sv *Server) {
	sv = &Server{
		entery: &fasthttp.Server{},
		router: fasthttprouter.New(),
		DB:     database.New(),
	}
	sv.entery.Handler = sv.router.Handler
	return sv
}

func (sv *Server) Start(port string) error {
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

//*****************| Utiles

func (sv *Server) StaticServer(ctx *fasthttp.RequestCtx, _ *Server) {
	_, body, _ := fasthttp.Get([]byte{}, fmt.Sprintf("http://localhost:3000%s", ctx.Path()))
	ctx.SetContentType("text/html")
	ctx.Write(body)
}
