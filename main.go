package main

import (
	"Wave/server"
	"Wave/server/api"
)

func main() {
	srv := server.New("./resources/configs/main.json")
	println("-- -- -- -- -- -- -- -- -- -- --")
	println("-- -- -- Server started -- -- --")
	println("-- -- -- -- -- -- -- -- -- -- --")

	api.UseAPI(srv)
	srv.Start()
}
