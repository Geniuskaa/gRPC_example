package server

import (
	proto "airTickets/pkg/gen/proto/v1"
	"airTickets/pkg/gen/proto/v1/airSale"
	"airTickets/pkg/gen/proto/v1/fewServs"
	"airTickets/pkg/gen/proto/v1/httpSample"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	Srv *grpc.Server

	airSaleSrv    airSale.Service
	fewServsSrv   fewServs.Service
	httpSampleSrv httpSample.Service
}

// Добавить логгер через StreamInterceptor
func NewGRPCServer(ctx context.Context, pools []*pgxpool.Pool, opt ...grpc.ServerOption) *gRPCServer {
	grpcSrv := grpc.NewServer()

	return &gRPCServer{Srv: grpcSrv,
		airSaleSrv:    *airSale.NewService(ctx, pools),
		fewServsSrv:   *fewServs.NewService(),
		httpSampleSrv: *httpSample.NewService(),
	}
}

func (g *gRPCServer) Init() {
	proto.RegisterHttpSampleServer(g.Srv, g.httpSampleSrv)
	proto.RegisterAirTicketsServiceServer(g.Srv, g.airSaleSrv)
	proto.RegisterMainServer(g.Srv, g.fewServsSrv)
	proto.RegisterMinorServer(g.Srv, g.fewServsSrv)
}

func (g *gRPCServer) GetSrv() *grpc.Server {
	return g.Srv
}
