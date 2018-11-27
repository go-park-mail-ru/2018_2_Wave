package grpcserver

import (
	"Wave/internal/config"
	lg "Wave/internal/logger"
	"google.golang.org/grpc"
	"net"
)

func StartServer(curlog *lg.Logger, GRPCC config.GRPCConfiguration) {
	lis, err := net.Listen("tcp", GRPCC.Port)
	if err != nil {

		curlog.Sugar.Infow("can't listen on port",
		"source", "grpcserv.go",
		"who", "New")
	}

	server := grpc.NewServer()
	//session.RegisterAuthCheckerServer(server, implementation.NewSessionManager(curlog))

	curlog.Sugar.Infow("starting grpc server on port" + GRPCC.Port,
		"source", "grpcserv.go",
		"who", "New")

	go server.Serve(lis)
}