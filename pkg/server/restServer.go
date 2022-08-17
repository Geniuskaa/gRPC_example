package server

import (
	proto "airTickets/pkg/gen/proto/v1"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"log"
	"net"
	"net/http"
)

type restServer struct {
	ctx context.Context
}

func NewRestServer(ctx context.Context, grpcSrv *gRPCServer) *http.Server {

	mux := runtime.NewServeMux()
	err := proto.RegisterHttpSampleHandlerServer(ctx, mux, grpcSrv.httpSampleSrv)

	if err != nil {
		log.Panic(err)
	}

	httpSrv := http.Server{
		Addr:    "localhost:8888",
		Handler: mux,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	return &httpSrv
}
