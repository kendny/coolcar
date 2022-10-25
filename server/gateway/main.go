package main

import (
	"context"
	authpb "coolcar/server/auth/api/gen/v1"
	rentalpb "coolcar/server/rental/api/gen/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
)

func main() {
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseEnumNumbers: true, // 将枚举类型生成数值， EnumsAsInts -> MarshalOptions.UseEnumNumbers
				UseProtoNames:  true, // 是否驼峰命名， OrigName -> MarshalOptions.UseProtoNames false：驼峰，true：下划线
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	))

	// 注册
	err := authpb.RegisterAuthServiceHandlerFromEndpoint(c,
		mux, // mux:multiplexer
		"localhost:8081",
		[]grpc.DialOption{grpc.WithInsecure()})

	if err != nil {
		log.Fatalf("cannot register auth service in grpc gateway:%v", err)
	}

	// 注册
	err = rentalpb.RegisterTripServiceHandlerFromEndpoint(c,
		mux, // mux:multiplexer
		"localhost:8082",
		[]grpc.DialOption{grpc.WithInsecure()})

	if err != nil {
		log.Fatalf("cannot register trip service in grpc gateway:%v", err)
	}

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("cannot listen and server:%v", err)
	}
}
