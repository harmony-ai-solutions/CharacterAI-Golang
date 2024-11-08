package cai

import (
	"github.com/harmony-ai-solutions/CharacterAI-Golang/cai"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

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

type ClientIntegrationSuite struct {
	suite.Suite
	client *cai.Client
	config *TestConfig
}

func TestClientIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ClientIntegrationSuite))
}

func (s *ClientIntegrationSuite) SetupSuite() {
	s.config = LoadTestConfig()

	if s.config.Token == "" || s.config.WebNextAuth == "" {
		s.T().Skip("API credentials are not provided. Skipping integration tests.")
	}

	s.client = cai.NewClient(s.config.Token, s.config.WebNextAuth, s.config.Proxy)

	// Authenticate the client
	err := s.client.Authenticate()
	s.Require().NoError(err, "Authentication failed")
	log.Debugf("Authenticated UserAccount for user ID '%v'", s.client.UserAccountID)

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)
}

func (s *ClientIntegrationSuite) TearDownSuite() {
	if s.client != nil {
		err := s.client.Close()
		s.Require().NoError(err, "Failed to close client")
	}
}

func (s *ClientIntegrationSuite) SetupTest() {
	// Set Log level for all tests
	log.SetLevel(log.DebugLevel)
}

func (s *ClientIntegrationSuite) TearDownTest() {
	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)
}
