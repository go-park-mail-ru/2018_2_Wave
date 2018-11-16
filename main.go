package main

import (
	"Wave/server"
)

//go:generate easyjson utiles/config/
//go:generate easyjson utiles/models/
//go:generate go run .

func main() {
	confPath := "./conf.json"
	server.Start(confPath)
}
