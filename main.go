package main

import (
	"Wave/server"
)

func main() {
	confPath := "./conf.json"
	srv := server.Init(confPath)

	server.Start(srv)
}
