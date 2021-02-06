package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HttpCreateTest struct {
	SuiteTest
}

func (s *HttpCreateTest) TestCreate() {
	tests := []struct {
		name  string
		event Event
	}{
		{
			"with notification",
			s.NewCommonEvent(),
		},
		{
			"without notification",
			func() Event {
				event := s.NewCommonEvent()
				event.Notification = nil
				return event
			}(),
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			data, _ := json.Marshal(tt.event)

			res, err := s.Call("create", data)
			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, res.StatusCode)
			id := s.readCreateId(res.Body)
			s.Require().Greater(id, 0)

			data, _ = json.Marshal(ListRequest{Date: tt.event.Start})

			res, err = s.Call("listday", data)
			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, res.StatusCode)
			events := s.readEvents(res.Body)
			s.Require().Equal(1, len(events))
			s.EqualEvents(tt.event, events[0])

			_ = s.app.DeleteAll(context.Background())
		})
	}
}

func (s *HttpCreateTest) TestCreateFail() {
	res, err := s.Call("create", []byte("Hello, world\n"))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, res.StatusCode)
}

func TestHttpCreateTest(t *testing.T) {
	suite.Run(t, new(HttpCreateTest))
}
