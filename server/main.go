/*
 * @Author: kendny wh_kendny@163.com
 * @Date: 2022-06-19 16:09:43
 * @LastEditors: kendny wh_kendny@163.com
 * @LastEditTime: 2022-06-21 13:46:26
 * @FilePath: /wx/Users/xxxian/go_project/src/coolcar/server/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"context"
	trippb "coolcar/server/proto/gen/go"
	trip "coolcar/server/tripservice"
	"log"
	"net"
	"net/http"

	// "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func startGRPCGateway() {
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseEnumNumbers: true,
				UseProtoNames:  true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	))
	err := trippb.RegisterTripServiceHandlerFromEndpoint(
		c,
		mux, // mux:multiplexer
		"localhost:8081",
		[]grpc.DialOption{grpc.WithInsecure()},
	)

	if err != nil {
		log.Fatalf("cannot start grpc gateway:%v", err)
	}

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("cannot listen and server:%v", err)
	}
}

func main() {
	// 设置一下log的形式
	log.SetFlags(log.Lshortfile)
	go startGRPCGateway()

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}

	s := grpc.NewServer()
	trippb.RegisterTripServiceServer(s, &trip.Service{})
	log.Fatal(s.Serve(lis))
}
