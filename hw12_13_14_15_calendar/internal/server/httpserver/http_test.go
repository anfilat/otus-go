package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
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
	ts   *httptest.Server
	app  app.App
	logg logger.Logger
	db   storage.Storage
}

func (s *SuiteTest) SetupTest() {
	ctx := context.Background()

	var buf bytes.Buffer
	s.logg, _ = logger.New("", &buf, "")

	dbConnect := os.Getenv("PQ_TEST")
	s.db, _ = initstorage.New(ctx, dbConnect == "", dbConnect)

	s.app = app.New(s.logg, s.db)

	s.ts = httptest.NewServer(newServer(s.app, s.logg).router)

	_ = s.app.DeleteAll(ctx)
}

func (s *SuiteTest) TearDownTest() {
	ctx := context.Background()
	s.ts.Close()
	_ = s.app.DeleteAll(ctx)
	_ = s.db.Close(ctx)
}

func (s *SuiteTest) Call(endPoint string, data []byte) (resp *http.Response, err error) {
	return http.Post(s.ts.URL+"/api/"+endPoint, "application/json", bytes.NewReader(data))
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
	s.Require().Equal(event1.UserID, event2.UserID)
	s.Require().Equal(event1.Notification, event2.Notification)
}

func (s *SuiteTest) AddEvent(event Event) int {
	data, _ := json.Marshal(event)

	res, err := http.Post(s.ts.URL+"/api/create", "application/json", bytes.NewReader(data))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	return s.readCreateId(res.Body)
}

func (s *SuiteTest) readCreateId(body io.ReadCloser) int {
	data, err := ioutil.ReadAll(body)
	defer body.Close()

	result := CreateResult{}
	err = json.Unmarshal(data, &result)
	s.Require().NoError(err)
	return result.ID
}

func (s *SuiteTest) readEvents(body io.ReadCloser) ListResult {
	data, err := ioutil.ReadAll(body)
	defer body.Close()

	result := ListResult{}
	err = json.Unmarshal(data, &result)
	s.Require().NoError(err)
	return result
}
