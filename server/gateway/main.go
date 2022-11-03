package main

import (
	"context"
	authpb "coolcar/server/auth/api/gen/v1"
	rentalpb "coolcar/server/rental/api/gen/v1"
	"coolcar/server/share/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
)

func main() {
	lg, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create zap logger: %v\n", err)
	}

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

	// grpc gateway 相对独立，不需要共用公共的RunGRPCServer
	serverConfig := []struct {
		name         string
		addr         string
		registerFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	}{
		{
			name:         "auth",
			addr:         "127.0.0:8081",
			registerFunc: authpb.RegisterAuthServiceHandlerFromEndpoint,
		},
		{
			name:         "rental",
			addr:         "127.0.0:8082",
			registerFunc: rentalpb.RegisterTripServiceHandlerFromEndpoint,
		},
	}

	for _, s := range serverConfig {
		err := s.registerFunc(c,
			mux, // mux:multiplexer
			"localhost:8081",
			[]grpc.DialOption{grpc.WithInsecure()})
		if err != nil {
			lg.Sugar().Fatalf("cannot register service %s:%v\n", s.name, s.addr)
		}
	}

	lg.Sugar().Fatal(http.ListenAndServe(":8080", mux))
}
