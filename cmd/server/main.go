package main

import (
	"airTickets/pkg/server"
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultPort = "9997"
	defaultHost = "0.0.0.0"

	countOfPool = 3
)

var errgRPCgracefullyStopped = errors.New("gRPC server was gracefully stopped.")

func main() {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGQUIT, syscall.SIGTERM)
	defer stop()

	//defaultDSNs := [...]string{"postgres://app:pass@localhost:5410/DB", "postgres://app:pass@localhost:5411/DB",
	//	"postgres://app:pass@localhost:5412/DB"}

	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	grpcSrv, gRPClistener, httpSrv := applicationStart(mainCtx, net.JoinHostPort(host, port), [3]string{"defaultDSNs", "", ""})

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		fmt.Println("HTTP server starting...")
		return httpSrv.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		fmt.Println("HTTP server is shutting down...")
		return httpSrv.Shutdown(context.Background())
	})
	g.Go(func() error {
		fmt.Println("gRPC server starting...")
		return grpcSrv.Serve(gRPClistener)
	})
	g.Go(func() error {
		<-gCtx.Done()
		fmt.Println("gRPC server is shutting down...")
		grpcSrv.GracefulStop()
		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
	fmt.Println("Servers were gracefully shut down.")
}

func applicationStart(ctx context.Context, addr string, pstgrs [countOfPool]string) (*grpc.Server, net.Listener, *http.Server) {
	gRPClistener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	////Создание пула подключений к PostgreSQL
	//poolArr := make([]*pgxpool.Pool, countOfPool)
	//for i := 0; i < countOfPool; i++ {
	//	poolArr[i], err = pgxpool.New(ctx, pstgrs[i])
	//	if err != nil {
	//		return err
	//	}
	//}

	//altsTC := alts.NewServerCreds(alts.DefaultServerOptions())

	grpcServer := server.NewGRPCServer(ctx, nil) // grpc.Creds(altsTC) , poolArr
	grpcServer.Init()

	restServer := server.NewRestServer(ctx, grpcServer)

	return grpcServer.GetSrv(), gRPClistener, restServer
}
