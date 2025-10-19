package main

import (
	"context"
	"net"
	"net/http"

	"github.com/bertguan96/cosf/api"
	"github.com/bertguan96/cosf/common"
	pb "github.com/bertguan96/cosf/cosf"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
)

func init() {
	common.InitLogger()
	common.InitDB()
}

func main() {
	grpclog.Info("server starting...")
	// 启动 GRPC 服务器
	go func() {
		lis, err := net.Listen("tcp", ":8000")
		if err != nil {
			grpclog.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterCosfServiceServer(grpcServer, &api.CosfServer{})
		if err := grpcServer.Serve(lis); err != nil {
			grpclog.Fatalf("failed to serve: %v", err)
		}
	}()

	// 启动 HTTP 服务器
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterCosfServiceHandlerFromEndpoint(ctx, mux, "localhost:8000", opts)
	if err != nil {
		panic(err)
	}

	// 启动 HTTP 服务器（如端口 8080）
	http.ListenAndServe(":8080", mux)
}
