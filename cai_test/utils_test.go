package cai

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type UtilsSuite struct {
	BaseSuite
}

func (s *UtilsSuite) TestPing() {
	reachable, err := s.client.Ping()
	s.Require().NoError(err, "Ping returned an error")
	s.Assert().True(reachable, "Service should be reachable")

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)
}

func TestUtilsSuite(t *testing.T) {
	suite.Run(t, new(UtilsSuite))
}
