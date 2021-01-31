package app_test

import (
	"bytes"
	"context"
	"os"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage/initstorage"
)

type SuiteTest struct {
	suite.Suite
	calendar app.App
	logg     logger.Logger
	db       storage.Storage
}

func (s *SuiteTest) SetupTest() {
	ctx := context.Background()

	var buf bytes.Buffer
	s.logg, _ = logger.New("", &buf, "")

	dbConnect := os.Getenv("PQ_TEST")
	s.db, _ = initstorage.New(ctx, dbConnect == "", dbConnect)

	s.calendar = app.New(s.logg, s.db)

	_ = s.calendar.DeleteAll(ctx)
}

func (s *SuiteTest) TearDownTest() {
	ctx := context.Background()
	_ = s.calendar.DeleteAll(ctx)
	_ = s.db.Close(ctx)
}

func (s *SuiteTest) NewCommonEvent() storage.Event {
	var eventStart = time.Now().Add(2 * time.Hour)
	var eventStop = eventStart.Add(time.Hour)
	notification := 4 * time.Hour

	return storage.Event{
		ID:           0,
		Title:        "some event",
		Start:        eventStart,
		Stop:         eventStop,
		Description:  "the event",
		UserID:       1,
		Notification: &notification,
	}
}

func (s *SuiteTest) AddEvent(event storage.Event) (int, error) {
	ctx := context.Background()
	id, err := s.calendar.Create(
		ctx,
		event.UserID,
		event.Title,
		event.Description,
		event.Start,
		event.Stop,
		event.Notification,
	)
	return id, err
}

func (s *SuiteTest) GetAll() []storage.Event {
	ctx := context.Background()
	data, err := s.calendar.ListAll(ctx)
	s.Require().NoError(err)
	return data
}

func (s *SuiteTest) EqualEvents(event1, event2 storage.Event) {
	s.Require().Equal(event1.Title, event2.Title)
	s.Require().Equal(event1.Description, event2.Description)
	s.Require().Equal(event1.Start.Unix(), event2.Start.Unix())
	s.Require().Equal(event1.Stop.Unix(), event2.Stop.Unix())
	s.Require().Equal(event1.UserID, event2.UserID)
	s.Require().Equal(event1.Notification, event2.Notification)
}
