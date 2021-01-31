package grpcserver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCUpdateTest struct {
	SuiteTest
}

func (s *GRPCUpdateTest) TestUpdate() {
	event := s.NewCommonEvent()
	id := s.AddEvent(event)

	event.Id = id
	event.Stop = timestamppb.New(event.Stop.AsTime().Add(time.Hour))

	ctx := context.Background()
	_, err := s.client.Update(ctx, event)
	s.Require().NoError(err)
}

func TestGRPCUpdateTest(t *testing.T) {
	suite.Run(t, new(GRPCUpdateTest))
}
