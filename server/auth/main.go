package main

import (
	"context"
	authpb "coolcar/server/auth/api/gen/v1"
	"coolcar/server/auth/auth"
	"coolcar/server/auth/dao"
	"coolcar/server/auth/wechat"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		logger.Fatal("cannot listen ", zap.Error(err))
	}

	// 数据库操作
	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://127.0.0.1:27017/?readPreference=primary&ssl=false&directConnection=true"))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}

	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, &auth.Service{
		OpenIDResolver: &wechat.Service{
			AppID:     "wxd5dda0926308f75e",
			AppSecret: "b986c1739cb5b8dcc6eb3c1badaf2735",
			Logger:    logger,
		},
		Mongo:  dao.NewMongo(mongoClient.Database("coolcar")),
		Logger: logger,
	})
	err = s.Serve(lis)
	logger.Fatal("cannot serve ", zap.Error(err))
}

func newZapLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = ""
	return cfg.Build()
}
