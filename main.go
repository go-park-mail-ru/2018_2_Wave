package main

import (
	"Wave/api"
	"Wave/server"
)

func main() {

	srv := server.New()

	api.UseAPI(srv)

	println("-- -- -- -- -- -- -- -- -- -- --")
	println("-- -- -- Server started -- -- --")
	println("-- -- -- -- -- -- -- -- -- -- --")

	srv.Start(":8080")
}
