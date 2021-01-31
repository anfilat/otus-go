package grpcserver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GRPCListTest struct {
	SuiteTest
}

func (s *GRPCListTest) TestListDay() {
	event := s.NewCommonEvent()
	s.AddEvent(event)

	ctx := context.Background()
	res, err := s.client.ListDay(ctx, &ListRequest{Date: event.Start})
	s.Require().NoError(err)
	s.Require().Equal(1, len(res.Events))
	s.EqualEvents(event, res.Events[0])
}

func (s *GRPCListTest) TestListWeek() {
	event := s.NewCommonEvent()
	s.AddEvent(event)

	ctx := context.Background()
	res, err := s.client.ListWeek(ctx, &ListRequest{Date: event.Start})
	s.Require().NoError(err)
	s.Require().Equal(1, len(res.Events))
	s.EqualEvents(event, res.Events[0])
}

func (s *GRPCListTest) TestListMonth() {
	event := s.NewCommonEvent()
	s.AddEvent(event)

	ctx := context.Background()
	res, err := s.client.ListMonth(ctx, &ListRequest{Date: event.Start})
	s.Require().NoError(err)
	s.Require().Equal(1, len(res.Events))
	s.EqualEvents(event, res.Events[0])
}

func TestGRPCListTest(t *testing.T) {
	suite.Run(t, new(GRPCListTest))
}
