package main

import (
	"Wave/server"
)

func main() {
	confPath := "./conf.json"
	server.Start(confPath)
	Wavelog.Sugar.Sync()
}
