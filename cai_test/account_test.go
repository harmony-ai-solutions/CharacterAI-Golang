package cai

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type AccountSuite struct {
	BaseSuite
}

func (s *AccountSuite) TestFetchMe() {
	account, err := s.client.FetchMe()
	s.Require().NoError(err, "FetchMe returned an error")
	s.Assert().NotNil(account, "Account should not be nil")
	s.Assert().NotEmpty(account.User.Username, "Username should not be empty")

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)
}

func TestAccountSuite(t *testing.T) {
	suite.Run(t, new(AccountSuite))
}
