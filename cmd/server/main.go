package main

import (
	"airTickets/cmd/server/app"
	airTickets "airTickets/pkg/airTicket/v1"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"net"
	"os"
	"time"
)

func main() {
	defaultDSNs := [...]string{"postgres://app:pass@localhost:5410/DB", "postgres://app:pass@localhost:5411/DB",
		"postgres://app:pass@localhost:5412/DB"}

	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	if err := execute(net.JoinHostPort(host, port), defaultDSNs); err != nil {
		os.Exit(1)
	}
}

const (
	defaultPort = "9997"
	defaultHost = "0.0.0.0"

	countOfPool = 3
)

func execute(addr string, pstgrs [countOfPool]string) (err error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	c := context.Background()
	ctx, cancel := context.WithCancel(c)

	go func() {
		time.Sleep(time.Second * 55)
		cancel()
	}()

	//Создание пула подключений к PostgreSQL
	poolArr := make([]*pgxpool.Pool, countOfPool)
	for i := 0; i < countOfPool; i++ {
		poolArr[i], err = pgxpool.Connect(ctx, pstgrs[i])
		if err != nil {
			return err
		}
	}

	grpcServer := grpc.NewServer()
	server := app.NewServer(poolArr)
	airTickets.RegisterAirTicketsServiceServer(grpcServer, server)

	return grpcServer.Serve(listener)
}
