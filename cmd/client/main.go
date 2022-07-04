package main

import (
	airTickets "airTickets/pkg/airTicket/v1"
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const defaultPort = "9997"
const defaultHost = "0.0.0.0"

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	if err := execute(net.JoinHostPort(host, port)); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(addr string) (err error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()

	client := airTickets.NewAirTicketsServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
	stream, err := client.AirTicketsFinder(ctx, &airTickets.TicketRequest{
		Date:                 &timestamp.Timestamp{Seconds: 1657016700},
		DepartureAirportCode: "KZN",
		ArrivalAirport:       "PEE",
	})

	if err != nil {
		return err
	}

	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		log.Print(response)
	}
	log.Print("finished")
	return nil

}
