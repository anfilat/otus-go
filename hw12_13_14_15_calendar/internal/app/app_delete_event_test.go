package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type DeleteEventTest struct {
	SuiteTest
}

func (s *DeleteEventTest) TestDeleteEvent() {
	event := s.NewCommonEvent()

	id1, err := s.AddEvent(event)
	s.Require().NoError(err)

	event.Start = event.Start.Add(2 * time.Hour)
	event.Stop = event.Start.Add(2 * time.Hour)
	id2, err := s.AddEvent(event)
	s.Require().NoError(err)

	ctx := context.Background()
	// удаление несуществующего события
	err = s.calendar.DeleteEvent(ctx, id2+1)
	s.Require().NoError(err)
	data, err := s.calendar.ListAllEvents(ctx)
	s.Require().NoError(err)
	s.Require().Equal(2, len(data))

	// удаление первого события
	err = s.calendar.DeleteEvent(ctx, id1)
	s.Require().NoError(err)
	data, err = s.calendar.ListAllEvents(ctx)
	s.Require().NoError(err)
	s.Require().Equal(1, len(data))

	// удаление второго события
	err = s.calendar.DeleteEvent(ctx, id2)
	s.Require().NoError(err)
	data, err = s.calendar.ListAllEvents(ctx)
	s.Require().NoError(err)
	s.Require().Equal(0, len(data))
}

func TestDeleteEventTest(t *testing.T) {
	suite.Run(t, new(DeleteEventTest))
}
