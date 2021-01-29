package httpserver

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HttpCreateTest struct {
	SuiteTest
}

func (s *HttpCreateTest) TestCreateEvent() {
	event := s.NewCommonEvent()
	data, err := json.Marshal(event)
	s.Require().NoError(err)

	res, err := http.Post(s.ts.URL+"/api/create", "application/json", bytes.NewReader(data))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	s.Require().Greater(s.readCreateId(res.Body), 0)

	data, err = json.Marshal(ListRequest{Date: event.Start})
	s.Require().NoError(err)

	res, err = http.Post(s.ts.URL+"/api/listday", "application/json", bytes.NewReader(data))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	events := s.readEvent(res.Body)
	s.Require().Equal(1, len(events))
	s.EqualEvents(event, events[0])
}

func (s *HttpCreateTest) readCreateId(body io.ReadCloser) int {
	data, err := ioutil.ReadAll(body)
	defer body.Close()

	result := CreateResult{}
	err = json.Unmarshal(data, &result)
	s.Require().NoError(err)
	return result.ID
}

func (s *HttpCreateTest) readEvent(body io.ReadCloser) ListResult {
	data, err := ioutil.ReadAll(body)
	defer body.Close()

	result := ListResult{}
	err = json.Unmarshal(data, &result)
	s.Require().NoError(err)
	return result
}

func TestHttpCreateTest(t *testing.T) {
	suite.Run(t, new(HttpCreateTest))
}
