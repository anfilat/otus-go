package grpcserver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GRPCDeleteTest struct {
	SuiteTest
}

func (s *GRPCDeleteTest) TestDelete() {
	event := s.NewCommonEvent()
	id := s.AddEvent(event)

	ctx := context.Background()
	_, err := s.client.Delete(ctx, &DeleteRequest{Id: id})
	s.Require().NoError(err)
}

func TestGRPCDeleteTest(t *testing.T) {
	suite.Run(t, new(GRPCDeleteTest))
}
