package httpSample

import (
	"airTickets/pkg/gen/proto/v1"
	"context"
)

type Service struct {
	v1.UnimplementedHttpSampleServer
}

func NewService() *Service {
	return &Service{}
}

func (s Service) StringResp(ctx context.Context, msg *v1.SimpleMsg) (*v1.SimpleMsg, error) {
	return &v1.SimpleMsg{
		Subject: msg.Subject,
		Body:    msg.Body,
	}, nil
}

func (s Service) StringGetReq(ctx context.Context, req *v1.Id) (*v1.SimpleMsgWithID, error) {
	return &v1.SimpleMsgWithID{
		Subject: "testdata",
		Body:    "testDataToo",
		Id:      req.Id,
	}, nil
}
