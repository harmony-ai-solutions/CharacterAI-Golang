package cai

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ChatSuite struct {
	BaseSuite
}

func (s *ChatSuite) TestCreateChatAndSendMessage() {
	// Create a new chat with a default character ID
	// Ensure the character ID is set
	if s.config.CharacterID == "" {
		s.T().Skip("No test character ID provided. Skipping TestCreateChatAndSendMessage.")
	}

	chat, greetingTurn, err := s.client.CreateChat(s.config.CharacterID, true)
	s.Require().NoError(err, "CreateChat returned an error")
	s.Assert().NotNil(chat, "Chat should not be nil")
	s.Assert().NotNil(greetingTurn, "Greeting turn should not be nil")

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)

	// Send a message
	messageText := "Hello, how are you?"
	turn, err := s.client.SendMessage(s.config.CharacterID, chat.ChatID, messageText)
	s.Require().NoError(err, "SendMessage returned an error")
	s.Assert().NotNil(turn, "Turn should not be nil")

	// Verify the response
	primaryCandidate := turn.Candidates[turn.PrimaryCandidateID]
	s.Require().NotNil(primaryCandidate, "Primary candidate should not be nil")
	s.Assert().NotEmpty(primaryCandidate.Text, "Response text should not be empty")

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)
}

func TestChatSuite(t *testing.T) {
	suite.Run(t, new(ChatSuite))
}
