package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type UpdateEventTest struct {
	SuiteTest
}

func (s *UpdateEventTest) TestUpdateEvent() {
	event := s.NewCommonEvent()
	id, err := s.AddEvent(event)
	s.Require().NoError(err)

	ctx := context.Background()
	updateEvent := storage.Event{
		Title:        "another event",
		Start:        time.Now().Add(5 * time.Hour),
		Stop:         time.Now().Add(6 * time.Hour),
		Description:  "very long event",
		UserID:       event.UserID,
		Notification: nil,
	}
	err = s.calendar.Update(ctx, id, updateEvent)
	s.Require().NoError(err)

	data := s.GetAll()
	s.Require().Equal(1, len(data))

	resultEvent := data[0]
	s.Require().Equal(id, resultEvent.ID)
	s.EqualEvents(updateEvent, resultEvent)
}

func (s *UpdateEventTest) TestUpdateEventFailNotExistsEvent() {
	event := s.NewCommonEvent()
	id, err := s.AddEvent(event)
	s.Require().NoError(err)

	event.Start = event.Start.Add(3 * time.Hour)
	event.Start = event.Stop.Add(3 * time.Hour)
	ctx := context.Background()
	err = s.calendar.Update(ctx, id+1, event)
	s.Require().Equal(storage.ErrNotExistsEvent, err)
}

func (s *UpdateEventTest) TestUpdateEventFailNoTitle() {
	event := s.NewCommonEvent()
	id, err := s.AddEvent(event)
	s.Require().NoError(err)

	event.Title = ""
	ctx := context.Background()
	err = s.calendar.Update(ctx, id, event)
	s.Require().Equal(app.ErrEmptyTitle, err)
}

func (s *UpdateEventTest) TestUpdateEventFailStartInPast() {
	event := s.NewCommonEvent()
	id, err := s.AddEvent(event)
	s.Require().NoError(err)

	event.Start = time.Now().Add(-time.Minute)
	ctx := context.Background()
	err = s.calendar.Update(ctx, id, event)
	s.Require().Equal(app.ErrStartInPast, err)
}

func (s *UpdateEventTest) TestUpdateEventNoDateBusy() {
	event := s.NewCommonEvent()
	id, err := s.AddEvent(event)
	s.Require().NoError(err)

	event2 := s.NewCommonEvent()
	event2.Start = event.Start.Add(2 * time.Hour)
	event2.Stop = event.Stop.Add(2 * time.Hour)
	_, err = s.AddEvent(event2)
	s.Require().NoError(err)

	ctx := context.Background()
	tests := []struct {
		start time.Time
		stop  time.Time
	}{
		{event.Start.Add(-30 * time.Minute), event.Start.Add(-20 * time.Minute)},
		{event.Start.Add(1 * time.Hour), event.Start.Add(2 * time.Hour)},
	}
	for _, tt := range tests {
		updateEvent := s.NewCommonEvent()
		updateEvent.Start = tt.start
		updateEvent.Stop = tt.stop
		err := s.calendar.Update(ctx, id, updateEvent)
		s.Require().NoError(err)
	}
}

func (s *UpdateEventTest) TestUpdateEventFailDateBusy() {
	event := s.NewCommonEvent()
	id, err := s.AddEvent(event)
	s.Require().NoError(err)

	event2 := s.NewCommonEvent()
	event2.Start = event2.Start.Add(3 * time.Hour)
	event2.Start = event2.Stop.Add(3 * time.Hour)
	_, err = s.AddEvent(event2)
	s.Require().NoError(err)

	ctx := context.Background()
	tests := []struct {
		start time.Time
		stop  time.Time
	}{
		{event2.Start.Add(-time.Hour), event2.Stop},
		{event2.Start.Add(-10 * time.Minute), event2.Start.Add(10 * time.Minute)},
		{event2.Stop.Add(-10 * time.Minute), event2.Stop.Add(10 * time.Minute)},
		{event2.Start.Add(10 * time.Minute), event2.Stop.Add(-10 * time.Minute)},
		{event2.Start.Add(-10 * time.Minute), event2.Stop.Add(10 * time.Minute)},
	}
	for _, tt := range tests {
		updateEvent := s.NewCommonEvent()
		updateEvent.Start = tt.start
		updateEvent.Stop = tt.stop
		err := s.calendar.Update(ctx, id, updateEvent)
		s.Require().Equal(app.ErrDateBusy, err)
	}
}

func TestUpdateEventTest(t *testing.T) {
	suite.Run(t, new(UpdateEventTest))
}
