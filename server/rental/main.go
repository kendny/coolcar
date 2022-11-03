package main

import (
	"coolcar/server/share/server"
	"log"
)

const publicKeyPath = "/Users/xxxian/go_project/src/coolcar/server/share/auth/public.key"

func main() {
	logger, err := server.NewZapLogger() //zap.NewDevelopment()

	if err != nil {
		log.Fatal("cannot create logger: %v", err)
	}

	// logger.Sugar() 带语法糖的logger
	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Addr:              ":8082",
		AuthPublicKeyFile: publicKeyPath,
		Logger:            logger,
		//RegisterFunc: func(s *grpc.Server) {
		//	rentalpb.RegisterTripServiceServer(s, &trip.Service{
		//		Logger: logger,
		//	})
		//},
	}))
}
