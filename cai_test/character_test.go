package cai

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CharacterSuite struct {
	BaseSuite
}

func (s *CharacterSuite) TestFetchCharacterInfo() {
	if s.config.CharacterID == "" {
		s.T().Skip("No test character ID provided. Skipping TestFetchCharacterInfo.")
	}

	character, err := s.client.FetchCharacterInfo(s.config.CharacterID)
	s.Require().NoError(err, "FetchCharacterInfo returned an error")
	s.Assert().NotNil(character, "Character should not be nil")
	s.Assert().Equal(s.config.CharacterID, character.CharacterID, "CharacterID should match")

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)
}

func TestCharacterSuite(t *testing.T) {
	suite.Run(t, new(CharacterSuite))
}
