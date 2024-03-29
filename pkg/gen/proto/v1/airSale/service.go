package airSale

import (
	"airTickets/pkg/gen/proto/v1"
	"context"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math/rand"
	"time"
)

type Service struct {
	ctx   context.Context
	pools []*pgxpool.Pool
	v1.UnimplementedAirTicketsServiceServer
}

func NewService(ctx context.Context, pools []*pgxpool.Pool) *Service {
	return &Service{ctx: ctx, pools: pools}
}

func (s Service) AirTicketsFinder(ticket *v1.TicketRequest, stream v1.AirTicketsService_AirTicketsFinderServer) error {
	log.Print(ticket)

	funcCtx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	//Как дать знать принимающей стороне, что ждать не стоит из-за возникшей ошибки? - Вернуть err тогда соедининение разорвется

	ch := make(chan v1.ProperFlightTicket, 3)

	for i := 0; i < 3; i++ {

		connCtx, _ := context.WithTimeout(funcCtx, time.Second*50)
		conn, err := s.pools[i].Acquire(connCtx)
		if err != nil {
			return status.Errorf(codes.Internal, "Problems with Pool connection")
		}

		go func(conn *pgxpool.Conn) {
			defer conn.Release()
			latency := rand.Int63n(10) + 1
			time.Sleep(time.Second * time.Duration(latency))
			row := conn.QueryRow(connCtx, `SELECT id,departure_time,flying_time,ticket_cost from air_tickets 
				where departure_airport=$1 and arrival_airport=$2 `, ticket.DepartureAirportCode, ticket.ArrivalAirport)

			var ticket = v1.ProperFlightTicket{
				Id:            0,
				DepartureTime: &timestamp.Timestamp{Seconds: 10},
				FlyingTime:    &duration.Duration{Seconds: 10},
				TicketCost:    0,
			}

			type ticketModel struct {
				Id            int64
				DepartureTime time.Time
				FlyingTime    time.Duration
				TicketCost    float64
			}
			var t ticketModel

			err := row.Scan(&t.Id, &t.DepartureTime, &t.FlyingTime, &t.TicketCost)
			if err != nil {
				log.Print(err)
				return
			}

			ticket.Id = t.Id
			ticket.DepartureTime.Seconds = t.DepartureTime.Unix()
			ticket.FlyingTime.Seconds = t.FlyingTime.Microseconds()
			ticket.TicketCost = t.TicketCost

			ch <- ticket
			log.Print("I am goroutine and i finish work")

			return

		}(conn)

	}

	lostPackets := 3
	for i := 0; i < 3; i++ {
		select {
		case t := <-ch:
			err := stream.Send(&t)
			if err != nil {
				continue
			}
		case <-time.After(time.Second * 25):
			lostPackets--
			continue
		}
	}

	if lostPackets < 1 {
		return status.Errorf(codes.NotFound, "Problems with such request")
	}

	return nil
}
