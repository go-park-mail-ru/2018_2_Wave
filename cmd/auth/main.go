package main

import (
	"Wave/internal/config"
	"Wave/internal/database"
	lg "Wave/internal/logger"
	au "Wave/internal/services/auth"
	"Wave/internal/services/auth/proto"

	"net"

	"google.golang.org/grpc"
)

const (
	confPath = "./configs/conf.json"

	logPath    = "./logs/"
	logFile = "auth-serv-log"
)


func main() {
	conf := config.Configure(confPath)
	curlog := lg.Construct(logPath, logFile)
	db := database.New(curlog)

	lis, err := net.Listen("tcp", conf.AC.Port)
	if err != nil {

		curlog.Sugar.Infow("can't listen on port",
		"source", "main.go",)

	}

	server := grpc.NewServer()
	auth.RegisterAuthServer(server, au.NewAuthManager(curlog, db))

	curlog.Sugar.Infow("starting grpc server on " + conf.AC.Host + conf.AC.Port,
		"source", "main.go",)

	server.Serve(lis)
}
