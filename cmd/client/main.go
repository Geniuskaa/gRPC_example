package main

import (
	proto "airTickets/pkg/gen/proto/v1"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	//err = airTicketsFinder(conn)
	//if err != nil {
	//	return err
	//}

	err = stringResp(conn)
	if err != nil {
		return err
	}

	err = stringReqWithID(conn)
	if err != nil {
		return err
	}

	return nil

}

func airTicketsFinder(conn *grpc.ClientConn) error {
	client := proto.NewAirTicketsServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
	stream, err := client.AirTicketsFinder(ctx, &proto.TicketRequest{
		Date:                 &timestamp.Timestamp{Seconds: 1657016700},
		DepartureAirportCode: "KZN",
		ArrivalAirport:       "PEE",
	})

	if err != nil {
		log.Print(err)
		return err
	}

	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Print(err)
			return err
		}
		log.Print(response)
	}
	log.Print("finished")

	return nil
}

func stringResp(conn *grpc.ClientConn) error {
	client := proto.NewHttpSampleClient(conn)
	ctx := context.Background()

	msg, err := client.StringResp(ctx, &proto.SimpleMsg{
		Subject: "Новый запрос",
		Body:    "Что-то там",
	})
	if err != nil {
		log.Print(err)
		return err
	}

	log.Print(fmt.Sprintf("We received back this letter: %s", msg))
	return nil
}

func stringReqWithID(conn *grpc.ClientConn) error {
	client := proto.NewHttpSampleClient(conn)
	ctx := context.Background()

	msg, err := client.StringGetReq(ctx, &proto.Id{
		Id: 1698,
	})
	if err != nil {
		log.Print(err)
		return err
	}

	log.Print(fmt.Sprintf("We received back this letter: %s", msg))
	return nil
}
