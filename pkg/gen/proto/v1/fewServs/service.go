package fewServs

import (
	fewServs "airTickets/pkg/gen/proto/v1"
)

type Service struct {
	fewServs.UnimplementedMainServer
	fewServs.UnimplementedMinorServer
}

func NewService() *Service {
	return &Service{}
}

// Если в одном proto файле несколько сервисов, как лучше их делить? Все заимплементить в одном сервере
// или для каждого сервиса использовать свой сервер?
