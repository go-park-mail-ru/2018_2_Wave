package main

import (
	"Wave/api"
	"Wave/server"
)

func main() {
	
	srv := server.New()
	
	api.UseAPI(srv)
	srv.GET("/index.html", srv.StaticServer)
	srv.GET("/app.bundle.js", srv.StaticServer)
	
	println("-- -- -- -- -- -- -- -- -- -- --")
	println("-- -- -- Server started -- -- --")
	println("-- -- -- -- -- -- -- -- -- -- --")

	srv.Start(":8080")
}
