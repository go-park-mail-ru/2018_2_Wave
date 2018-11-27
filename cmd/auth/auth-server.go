package main

import (
	"Wave/internal/grpcserver"
	"Wave/internal/config"
	lg "Wave/internal/logger"
)

func main() {
	path := "./configs/conf.json"
	conf := config.Configure(path)
	curlog := lg.Construct()
	prof := mc.Construct()
	db := database.New(curlog)

	grpcserver.StartServer(curlog, conf.GRPCC, &db)
}
