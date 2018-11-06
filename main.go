package main

import (
	"Wave/server"
	"Wave/utiles/logger"
)

func main() {
	confPath := "./conf.json"
	wavelog := logger.Construct()
	server.Start(confPath, wavelog)
	wavelog.Sugar.Sync()
}
