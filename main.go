package main

import (
	"test_module/server"
)

func main() {
	srv := server.Server{}
	srv.Init()
	srv.Get("/gg", srv.StaticServer)
	srv.Get("/wp/:path", srv.StaticServer)
	srv.Start(":8080")
}
