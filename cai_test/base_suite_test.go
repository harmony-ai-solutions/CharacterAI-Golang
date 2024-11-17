package cai

import (
	"github.com/harmony-ai-solutions/CharacterAI-Golang/cai"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"os"
	"time"
)

type BaseSuite struct {
	suite.Suite
	client *cai.Client
	config *TestConfig
}

type TestConfig struct {
	Token       string
	WebNextAuth string
	Proxy       string
	// Character Testing
	CharacterID string
}

func LoadTestConfig() *TestConfig {
	return &TestConfig{
		Token:       os.Getenv("CHARACTERAI_TOKEN"),
		WebNextAuth: os.Getenv("CHARACTERAI_WEBNEXTAUTH"),
		Proxy:       os.Getenv("CHARACTERAI_PROXY"),
		CharacterID: os.Getenv("CHARACTERAI_CHARACTERID"),
	}
}

func (s *BaseSuite) SetupSuite() {
	s.config = LoadTestConfig()
	s.client = cai.NewClient(s.config.Token, s.config.WebNextAuth, s.config.Proxy)
	err := s.client.Authenticate()
	s.Require().NoError(err)
}

func (s *BaseSuite) TearDownSuite() {
	if s.client != nil {
		err := s.client.Close()
		s.Require().NoError(err, "Failed to close client")
	}
}

func (s *BaseSuite) SetupTest() {
	// Set Log level for all tests
	log.SetLevel(log.DebugLevel)
}

func (s *BaseSuite) TearDownTest() {
	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)
}
