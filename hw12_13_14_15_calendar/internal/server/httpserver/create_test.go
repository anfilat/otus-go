package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HttpCreateTest struct {
	SuiteTest
}

func (s *HttpCreateTest) TestCreate() {
	event := s.NewCommonEvent()
	data, _ := json.Marshal(event)

	res, err := http.Post(s.ts.URL+"/api/create", "application/json", bytes.NewReader(data))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	s.Require().Greater(s.readCreateId(res.Body), 0)

	data, _ = json.Marshal(ListRequest{Date: event.Start})

	res, err = http.Post(s.ts.URL+"/api/listday", "application/json", bytes.NewReader(data))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	events := s.readEvents(res.Body)
	s.Require().Equal(1, len(events))
	s.EqualEvents(event, events[0])
}

func (s *HttpCreateTest) TestCreateFail() {
	res, err := http.Post(s.ts.URL+"/api/create", "application/json", bytes.NewReader([]byte("Hello, world\n")))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, res.StatusCode)
}

func TestHttpCreateTest(t *testing.T) {
	suite.Run(t, new(HttpCreateTest))
}
