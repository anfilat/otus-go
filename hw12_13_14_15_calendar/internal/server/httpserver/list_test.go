package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HttpListTest struct {
	SuiteTest
}

func (s *HttpListTest) TestListDay() {
	event := s.NewCommonEvent()
	s.AddEvent(event)

	data, _ := json.Marshal(ListRequest{Date: event.Start})

	res, err := http.Post(s.ts.URL+"/api/listday", "application/json", bytes.NewReader(data))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	events := s.readEvents(res.Body)
	s.Require().Equal(1, len(events))
	s.EqualEvents(event, events[0])
}

func (s *HttpListTest) TestListWeek() {
	event := s.NewCommonEvent()
	s.AddEvent(event)

	data, _ := json.Marshal(ListRequest{Date: event.Start})

	res, err := http.Post(s.ts.URL+"/api/listweek", "application/json", bytes.NewReader(data))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	events := s.readEvents(res.Body)
	s.Require().Equal(1, len(events))
	s.EqualEvents(event, events[0])
}

func (s *HttpListTest) TestListMonth() {
	event := s.NewCommonEvent()
	s.AddEvent(event)

	data, _ := json.Marshal(ListRequest{Date: event.Start})

	res, err := http.Post(s.ts.URL+"/api/listmonth", "application/json", bytes.NewReader(data))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	events := s.readEvents(res.Body)
	s.Require().Equal(1, len(events))
	s.EqualEvents(event, events[0])
}

func TestHttpListTest(t *testing.T) {
	suite.Run(t, new(HttpListTest))
}
