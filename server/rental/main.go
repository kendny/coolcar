package main

import (
	"context"
	rentalpb "coolcar/server/rental/api/gen/v1"
	"coolcar/server/rental/trip"
	"coolcar/server/rental/trip/client/car"
	"coolcar/server/rental/trip/client/poi"
	"coolcar/server/rental/trip/client/profile"
	"coolcar/server/rental/trip/dao"
	"coolcar/server/share/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
)

const publicKeyPath = "/Users/xxxian/go_project/src/coolcar/server/share/auth/public.key"

func main() {
	logger, err := server.NewZapLogger() //zap.NewDevelopment()

	if err != nil {
		log.Fatal("cannot create logger: %v", err)
	}

	// 数据库操作
	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://127.0.0.1:27017/?readPreference=primary&ssl=false&directConnection=true"))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}

	// logger.Sugar() 带语法糖的logger
	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Addr:              ":8082",
		AuthPublicKeyFile: publicKeyPath,
		Logger:            logger,
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				CarManager:     &car.Manager{},
				ProfileManager: &profile.Manager{},
				POIManager:     &poi.Manager{},
				Mongo:          dao.NewMongo(mongoClient.Database("coolcar")),
				Logger:         logger,
			})
		},
	}))
}
