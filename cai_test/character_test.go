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
	s.Assert().Equal(s.config.CharacterID, character.ExternalID, "ExternalID should match")

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)
}

// TestFetchCharactersByCategory tests the FetchCharactersByCategory method
func (s *CharacterSuite) TestFetchCharactersByCategory() {
	categories, err := s.client.FetchCharactersByCategory()
	s.Require().NoError(err, "FetchCharactersByCategory returned an error")
	s.Assert().NotNil(categories, "Categories should not be nil")
	s.Assert().NotEmpty(categories, "Categories should not be empty")

	// Check that each category has at least one character
	for category, characters := range categories {
		s.Assert().NotEmpty(characters, "Category '%s' should have characters", category)
	}

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)
}

// TestFetchRecommendedCharacters tests the FetchRecommendedCharacters method
func (s *CharacterSuite) TestFetchRecommendedCharacters() {
	characters, err := s.client.FetchRecommendedCharacters()
	s.Require().NoError(err, "FetchRecommendedCharacters returned an error")
	s.Assert().NotNil(characters, "Characters should not be nil")
	s.Assert().NotEmpty(characters, "Characters should not be empty")

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)
}

// TestFetchFeaturedCharacters tests the FetchFeaturedCharacters method
func (s *CharacterSuite) TestFetchFeaturedCharacters() {
	characters, err := s.client.FetchFeaturedCharacters()
	s.Require().NoError(err, "FetchFeaturedCharacters returned an error")
	s.Assert().NotNil(characters, "Characters should not be nil")
	s.Assert().NotEmpty(characters, "Characters should not be empty")

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)
}

// TestFetchSimilarCharacters tests the FetchSimilarCharacters method
func (s *CharacterSuite) TestFetchSimilarCharacters() {
	if s.config.CharacterID == "" {
		s.T().Skip("No test character ID provided. Skipping TestFetchSimilarCharacters.")
	}

	characters, err := s.client.FetchSimilarCharacters(s.config.CharacterID)
	s.Require().NoError(err, "FetchSimilarCharacters returned an error")
	s.Assert().NotNil(characters, "Characters should not be nil")
	//s.Assert().NotEmpty(characters, "Characters should not be empty")

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)
}

// TestSearchCharacters tests the SearchCharacters method
func (s *CharacterSuite) TestSearchCharacters() {
	query := "test"
	characters, err := s.client.SearchCharacters(query)
	s.Require().NoError(err, "SearchCharacters returned an error")
	s.Assert().NotNil(characters, "Characters should not be nil")
	s.Assert().NotEmpty(characters, "Characters should not be empty")

	// Check that the characters match the query
	//for _, character := range characters {
	//	s.Assert().Contains(character.Name, query, "Character name should contain the query")
	//}

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)
}

func TestCharacterSuite(t *testing.T) {
	suite.Run(t, new(CharacterSuite))
}
