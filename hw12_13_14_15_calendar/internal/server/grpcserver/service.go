package grpcserver

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Service struct {
	UnimplementedCalendarServer

	app app.App
}

func NewService(app app.App) *Service {
	return &Service{
		app: app,
	}
}

func (s *Service) Create(ctx context.Context, req *Event) (*CreateResult, error) {
	id, err := s.app.Create(
		ctx,
		int(req.UserId),
		req.Title,
		req.Description,
		req.Start.AsTime(),
		req.Stop.AsTime(),
		getNotification(req),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &CreateResult{Id: int32(id)}, nil
}

func (s *Service) Update(ctx context.Context, req *Event) (*UpdateResult, error) {
	change := storage.Event{
		ID:           int(req.Id),
		Title:        req.Title,
		Start:        req.Start.AsTime(),
		Stop:         req.Stop.AsTime(),
		Description:  req.Description,
		UserID:       int(req.UserId),
		Notification: getNotification(req),
	}
	err := s.app.Update(ctx, int(req.Id), change)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &UpdateResult{}, nil
}

func getNotification(req *Event) *time.Duration {
	if req.Notification != nil {
		data := req.Notification.AsDuration()
		return &data
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResult, error) {
	err := s.app.Delete(ctx, int(req.Id))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &DeleteResult{}, nil
}

func (s *Service) ListDay(ctx context.Context, req *ListRequest) (*ListResult, error) {
	events, err := s.app.ListDay(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ListResult{Events: storageEventsToGRPCEvents(events)}, nil
}

func (s *Service) ListWeek(ctx context.Context, req *ListRequest) (*ListResult, error) {
	events, err := s.app.ListWeek(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ListResult{Events: storageEventsToGRPCEvents(events)}, nil
}

func (s *Service) ListMonth(ctx context.Context, req *ListRequest) (*ListResult, error) {
	events, err := s.app.ListMonth(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ListResult{Events: storageEventsToGRPCEvents(events)}, nil
}

func storageEventsToGRPCEvents(events []storage.Event) []*Event {
	resultEvents := make([]*Event, 0, len(events))
	for _, event := range events {
		resultEvent := &Event{
			Id:          int32(event.ID),
			Title:       event.Title,
			Start:       timestamppb.New(event.Start),
			Stop:        timestamppb.New(event.Stop),
			Description: event.Description,
			UserId:      int32(event.UserID),
		}
		if event.Notification != nil {
			notification := *event.Notification
			resultEvent.Notification = durationpb.New(notification)
		}
		resultEvents = append(resultEvents, resultEvent)
	}
	return resultEvents
}
