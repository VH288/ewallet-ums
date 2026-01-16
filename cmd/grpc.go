package cmd

import (
	"log"
	"net"

	"ewallet-ums/cmd/proto/tokenvalidation"
	"ewallet-ums/helpers"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func ServeGRPC() {
	// init dependency
	dependency := dependencyInject()

	lis, err := net.Listen("tcp", ":"+helpers.GetEnv("GRPC_PORT", "7000"))
	if err != nil {
		log.Fatal("failed to listen grpc: ", err)
	}

	s := grpc.NewServer()

	// list method
	tokenvalidation.RegisterTokenValidationServer(s, dependency.TokenValidationAPI)

	logrus.Info("start listening grpc on port: " + helpers.GetEnv("GRPC_PORT", "7000"))
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve grpc port: ", err)
	}
}
