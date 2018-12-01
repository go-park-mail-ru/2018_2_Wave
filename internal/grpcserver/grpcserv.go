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