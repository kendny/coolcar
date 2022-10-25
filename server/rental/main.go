package main

import (
	rentalpb "coolcar/server/rental/api/gen/v1"
	"coolcar/server/rental/trip"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	logger, err := newZapLogger() //zap.NewDevelopment()

	if err != nil {
		log.Fatal("cannot create logger: %v", err)
	}

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		logger.Fatal("cannot listen ", zap.Error(err))
	}

	s := grpc.NewServer()
	rentalpb.RegisterTripServiceServer(s, &trip.Service{
		Logger: logger,
	})

	err = s.Serve(lis)
	logger.Fatal("cannot serve ", zap.Error(err))
}

// 配置日志
func newZapLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = ""
	return cfg.Build()
}
