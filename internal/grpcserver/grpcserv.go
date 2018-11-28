package grpcserver

import (
	"Wave/internal/config"
	"Wave/internal/database"
	lg "Wave/internal/logger"
	mc "Wave/internal/metrics"
	au "Wave/internal/services/auth"
	gm "Wave/internal/services/game"
	"Wave/internal/services/game/proto"
	"Wave/internal/services/auth/proto"

	"net"

	"google.golang.org/grpc"
)

func startAuthServer(curlog *lg.Logger, GRPCC config.GRPCConfiguration, db *database.DatabaseModel) {
	lis, err := net.Listen("tcp", GRPCC.Port)
	if err != nil {

		curlog.Sugar.Infow("can't listen on port",
		"source", "grpcserv.go",
		"who", "New")

	}

	server := grpc.NewServer()
	auth.RegisterAuthServer(server, au.NewAuthManager(curlog, db))

	curlog.Sugar.Infow("starting grpc server on port" + GRPCC.Port,
		"source", "grpcserv.go",
		"who", "New")

	go server.Serve(lis)
}

func startGameServer(curlog *lg.Logger, conf config.Configuration, Prof *mc.Profiler) {
	lis, err := net.Listen("tcp", conf.Game.Port)
	if err != nil {
		curlog.Sugar.Infow("can't listen on port",
		"source", "grpcserv.go",
		"who", "New")
	}
	server := grpc.NewServer()
	game.RegisterGameServer(server, gm.NewGame(curlog, Prof, conf))
	go server.Serve(lis)
}

// StartServer - start microservices
func StartServer(curlog *lg.Logger, conf config.Configuration, db *database.DatabaseModel, Prof *mc.Profiler) {
	startAuthServer(curlog, conf.Auth, db)
	startGameServer(curlog, conf, Prof)
}
