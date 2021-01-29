package httpserver

import (
	"bytes"
	"context"
	"net/http/httptest"
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
	db storage.Storage
	ts *httptest.Server
}

func (s *SuiteTest) SetupTest() {
	var buf bytes.Buffer
	logg, _ := logger.New("", &buf, "")

	ctx := context.Background()
	dbConnect := os.Getenv("PQ_TEST")
	db, _ := initstorage.New(ctx, dbConnect == "", dbConnect)
	s.db = db

	_ = s.db.DeleteAll(ctx)

	calendar := app.New(logg, db)
	s.ts = httptest.NewServer(newServer(calendar, logg).router)
}

func (s *SuiteTest) TearDownTest() {
	ctx := context.Background()
	s.ts.Close()
	_ = s.db.DeleteAll(ctx)
	_ = s.db.Close(ctx)
}

func (s *SuiteTest) NewCommonEvent() Event {
	var eventStart = time.Now().Add(2 * time.Hour)
	var eventStop = eventStart.Add(time.Hour)
	notification := 4 * time.Hour

	return Event{
		ID:           0,
		Title:        "some event",
		Start:        eventStart,
		Stop:         eventStop,
		Description:  "the event",
		UserID:       1,
		Notification: &notification,
	}
}

func (s *SuiteTest) EqualEvents(event1, event2 Event) {
	s.Require().Equal(event1.Title, event2.Title)
	s.Require().Equal(event1.Description, event2.Description)
	s.Require().Equal(event1.Start.Unix(), event2.Start.Unix())
	s.Require().Equal(event1.Stop.Unix(), event2.Stop.Unix())
	s.Require().Equal(event1.Notification, event2.Notification)
}
