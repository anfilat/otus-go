package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type HttpUpdateTest struct {
	SuiteTest
}

func (s *HttpUpdateTest) TestUpdate() {
	event := s.NewCommonEvent()
	id := s.AddEvent(event)

	event.ID = id
	event.Stop = event.Stop.Add(time.Hour)
	data, _ := json.Marshal(event)

	res, err := http.Post(s.ts.URL+"/api/update", "application/json", bytes.NewReader(data))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
}

func TestHttpUpdateTest(t *testing.T) {
	suite.Run(t, new(HttpUpdateTest))
}
