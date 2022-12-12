package main

import (
	"context"
	blobpb "coolcar/server/blob/api/gen/v1"
	"coolcar/server/blob/blob"
	"coolcar/server/blob/cos"

	"coolcar/server/blob/dao"
	"coolcar/server/share/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
)

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		logger.Fatal("cannot connect to mongodb", zap.Error(err))
	}
	db := mongoClient.Database("coolcar")
	st, err := cos.NewService(
		"https://wuhan-1259722894.cos.ap-shanghai.myqcloud.com",
		"AKIDbkfNr78vUq32pOhoiQxHMDpDPPESeicR",
		"jHiz62oV8lK1Zv78yeEGE90hHl48zc1B",
	)
	if err != nil {
		logger.Fatal("cannot create cos service", zap.Error(err))
	}
	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:   "blob",
		Addr:   ":8083",
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			blobpb.RegisterBlobServiceServer(s, &blob.Service{
				Storage: st,
				Mongo:   dao.NewMongo(db),
				Logger:  logger,
			})
		},
	}))
}
