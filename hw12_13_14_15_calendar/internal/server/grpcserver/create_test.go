package grpcserver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GRPCCreateTest struct {
	SuiteTest
}

func (s *GRPCCreateTest) TestCreate() {
	tests := []struct {
		name  string
		event *Event
	}{
		{
			"with notification",
			s.NewCommonEvent(),
		},
		{
			"without notification",
			func() *Event {
				event := s.NewCommonEvent()
				event.Notification = nil
				return event
			}(),
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			createRes, err := s.client.Create(ctx, tt.event)
			s.Require().NoError(err)
			s.Require().Greater(createRes.Id, int32(0))

			listRes, err := s.client.ListDay(ctx, &ListRequest{Date: tt.event.Start})
			s.Require().NoError(err)
			s.Require().Equal(1, len(listRes.Events))
			s.EqualEvents(tt.event, listRes.Events[0])

			_ = s.app.DeleteAll(context.Background())
		})
	}
}

func TestGRPCCreateTest(t *testing.T) {
	suite.Run(t, new(GRPCCreateTest))
}
