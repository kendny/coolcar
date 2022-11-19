package main

import (
	"context"
	rentalpb "coolcar/server/rental/api/gen/v1"
	"coolcar/server/rental/profile"
	profiledao "coolcar/server/rental/profile/dao"
	"coolcar/server/rental/trip"
	"coolcar/server/rental/trip/client/car"
	"coolcar/server/rental/trip/client/poi"
	profClient "coolcar/server/rental/trip/client/profile"
	tripdao "coolcar/server/rental/trip/dao"
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

	db := mongoClient.Database("coolcar")
	// logger.Sugar() 带语法糖的logger
	profService := &profile.Service{
		Mongo:  profiledao.NewMongo(db),
		Logger: logger,
	}
	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Addr:              ":8082",
		AuthPublicKeyFile: publicKeyPath,
		Logger:            logger,
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				CarManager: &car.Manager{},
				ProfileManager: &profClient.Manager{
					Fetcher: profService,
				},
				POIManager: &poi.Manager{},
				Mongo:      tripdao.NewMongo(db),
				Logger:     logger,
			})
			rentalpb.RegisterProfileServiceServer(s, profService)
		},
	}))
}
