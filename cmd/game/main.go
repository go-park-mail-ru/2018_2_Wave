package main

import (
	"Wave/internal/config"
	lg "Wave/internal/logger"
	mc "Wave/internal/metrics"
	gm "Wave/internal/services/game"
	"Wave/internal/services/game/proto"

	"net"

	"google.golang.org/grpc"
)

const (
	confPath = "./configs/conf.json"

	logPath    = "./logs/"
	logFile = "game-serv-log"
)

func main() {
	conf := config.Configure(confPath)
	curlog := lg.Construct(logPath, logFile)
	prof := mc.Construct()

	lis, err := net.Listen("tcp", conf.Game.Port)
	if err != nil {

		curlog.Sugar.Infow("can't listen on port",
		"source", "main.go",)

	}

	server := grpc.NewServer()
	game.RegisterGameServer(server, gm.NewGame(curlog, prof, conf))

	server.Serve(lis)
}