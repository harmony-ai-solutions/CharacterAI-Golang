package cai

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"os"
	"strconv"
	"testing"
)

type CAIClientTestSuite struct {
	suite.Suite
}

func TestCAIClientTestSuite(t *testing.T) {
	suite.Run(t, new(CAIClientTestSuite))
}

func (s *CAIClientTestSuite) SetupTest() {
	// Set Log level for all tests
	log.SetLevel(log.DebugLevel)
}

func (s *CAIClientTestSuite) TestCAILoginWrongCredentials() {
	// Init params
	token := os.Getenv("CAI_TOKEN")[1:4]
	isPlus, errParse := strconv.ParseBool(os.Getenv("CAI_PLUS"))
	if errParse != nil {
		isPlus = false
	}

	// Create client Wrapper
	cai, errClient := NewGoCAI(token, isPlus)
	s.Nil(errClient)
	s.NotNil(cai)

	// Perform simple Backend call
	userCharacters, errCharacters := cai.User.Characters()
	s.Nil(errCharacters)
	details, isDetailMessage := userCharacters["detail"]
	s.True(isDetailMessage)
	s.Equal("Authentication credentials were not provided.", details)
}

func (s *CAIClientTestSuite) TestCAILoginCorrectCredentials() {
	// Init params
	token := os.Getenv("CAI_TOKEN")
	isPlus, errParse := strconv.ParseBool(os.Getenv("CAI_PLUS"))
	if errParse != nil {
		isPlus = false
	}

	// Create client
	cai, errClient := NewGoCAI(token, isPlus)
	s.Nil(errClient)
	s.NotNil(cai)

	// Perform simple Backend call
	userCharacters, errCharacters := cai.User.Characters()
	s.Nil(errCharacters)
	characters, isCharacters := userCharacters["characters"]
	s.True(isCharacters)
	s.NotNil(characters)
}
