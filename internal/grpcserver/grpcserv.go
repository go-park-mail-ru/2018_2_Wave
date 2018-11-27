package grpcserver

import (
	"Wave/internal/config"
	"Wave/internal/database"
	lg "Wave/internal/logger"
	au "Wave/internal/services/auth"
	"Wave/internal/services/auth/proto"

	"net"

	"google.golang.org/grpc"
)

func StartServer(curlog *lg.Logger, GRPCC config.GRPCConfiguration, db *database.DatabaseModel) {
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