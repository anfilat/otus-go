package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/app"
)

type CreateEventTest struct {
	SuiteTest
}

func (s *CreateEventTest) TestCreateEvent() {
	event := s.NewCommonEvent()

	id, err := s.AddEvent(event)
	s.Require().NoError(err)
	s.Require().Greater(id, 0)

	data := s.GetAll()
	s.Require().Equal(1, len(data))

	resultEvent := data[0]
	s.EqualEvents(event, resultEvent)
}

func (s *CreateEventTest) TestCreateEventFailNoUser() {
	event := s.NewCommonEvent()

	event.UserID = 0
	_, err := s.AddEvent(event)
	s.Require().Equal(app.ErrNoUserID, err)
}

func (s *CreateEventTest) TestCreateEventFailNoTitle() {
	event := s.NewCommonEvent()

	event.Title = ""
	_, err := s.AddEvent(event)
	s.Require().Equal(app.ErrEmptyTitle, err)
}

func (s *CreateEventTest) TestCreateEventFailStartInPast() {
	event := s.NewCommonEvent()

	event.Start = time.Now().Add(-time.Minute)
	_, err := s.AddEvent(event)
	s.Require().Equal(app.ErrStartInPast, err)
}

func (s *CreateEventTest) TestCreateEventForOtherUser() {
	event := s.NewCommonEvent()
	_, err := s.AddEvent(event)
	s.Require().NoError(err)

	event.UserID = event.UserID + 1
	_, err = s.AddEvent(event)
	s.Require().NoError(err)
}

func (s *CreateEventTest) TestCreateEventNoDateBusy() {
	event := s.NewCommonEvent()
	_, err := s.AddEvent(event)
	s.Require().NoError(err)

	tests := []struct {
		start time.Time
		stop  time.Time
	}{
		{event.Start.Add(-30 * time.Minute), event.Start.Add(-20 * time.Minute)},
		{event.Start.Add(1 * time.Hour), event.Start.Add(2 * time.Hour)},
	}
	for _, tt := range tests {
		err := s.AddEventForTime(tt.start, tt.stop)
		s.Require().NoError(err)
	}
}

func (s *CreateEventTest) TestCreateEventFailDateBusy() {
	event := s.NewCommonEvent()
	_, err := s.AddEvent(event)
	s.Require().NoError(err)

	tests := []struct {
		start time.Time
		stop  time.Time
	}{
		{event.Start.Add(-time.Hour), event.Stop},
		{event.Start.Add(-10 * time.Minute), event.Start.Add(10 * time.Minute)},
		{event.Stop.Add(-10 * time.Minute), event.Stop.Add(10 * time.Minute)},
		{event.Start.Add(10 * time.Minute), event.Stop.Add(-10 * time.Minute)},
		{event.Start.Add(-10 * time.Minute), event.Stop.Add(10 * time.Minute)},
	}
	for _, tt := range tests {
		err := s.AddEventForTime(tt.start, tt.stop)
		s.Require().Equal(app.ErrDateBusy, err)
	}
}

func (s *CreateEventTest) AddEventForTime(start, stop time.Time) error {
	event := s.NewCommonEvent()
	ctx := context.Background()
	_, err := s.calendar.Create(
		ctx,
		event.UserID,
		event.Title,
		event.Description,
		start,
		stop,
		event.Notification,
	)
	return err
}

func TestCreateEventTest(t *testing.T) {
	suite.Run(t, new(CreateEventTest))
}
