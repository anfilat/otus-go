package httpserver

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HttpDeleteTest struct {
	SuiteTest
}

func (s *HttpDeleteTest) TestDelete() {
	event := s.NewCommonEvent()
	id := s.AddEvent(event)

	data, _ := json.Marshal(DeleteRequest{ID: id})

	res, err := s.Call("delete", data)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
}

func TestHttpDeleteTest(t *testing.T) {
	suite.Run(t, new(HttpDeleteTest))
}
