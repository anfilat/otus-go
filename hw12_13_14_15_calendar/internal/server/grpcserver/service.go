package grpcserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	UnimplementedCalendarServer
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Create(context.Context, *Event) (*EventCreateResult, error) {
	return nil, status.Error(codes.InvalidArgument, "oh, all wrong")
}
