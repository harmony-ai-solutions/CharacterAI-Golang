package cai

import (
	"github.com/harmony-ai-solutions/CharacterAI-Golang/cai"
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

func (s *ChatSuite) TestFetchHistories() {
	if s.config.CharacterID == "" {
		s.T().Skip("No test character ID provided. Skipping TestFetchHistories.")
	}

	histories, err := s.client.FetchHistories(s.config.CharacterID, 5)
	s.Require().NoError(err, "FetchHistories returned an error")
	s.Assert().NotNil(histories, "Histories should not be nil")
}

func (s *ChatSuite) TestFetchChats() {
	if s.config.CharacterID == "" {
		s.T().Skip("No test character ID provided. Skipping TestFetchChats.")
	}

	chats, err := s.client.FetchChats(s.config.CharacterID, 5)
	s.Require().NoError(err, "FetchChats returned an error")
	s.Assert().NotNil(chats, "Chats should not be nil")
}

func (s *ChatSuite) TestFetchChat() {
	if s.config.CharacterID == "" {
		s.T().Skip("No test character ID provided. Skipping TestFetchChat.")
	}

	// Create a new chat
	chat, _, err := s.client.CreateChat(s.config.CharacterID, true)
	s.Require().NoError(err, "CreateChat returned an error")
	s.Assert().NotNil(chat, "Chat should not be nil")

	// Fetch the chat
	fetchedChat, err := s.client.FetchChat(chat.ChatID)
	s.Require().NoError(err, "FetchChat returned an error")
	s.Assert().NotNil(fetchedChat, "Fetched chat should not be nil")
	s.Assert().Equal(chat.ChatID, fetchedChat.ChatID, "Fetched chat ID should match")
}

func (s *ChatSuite) TestFetchRecentChats() {
	chats, err := s.client.FetchRecentChats()
	s.Require().NoError(err, "FetchRecentChats returned an error")
	s.Assert().NotNil(chats, "Chats should not be nil")
}

func (s *ChatSuite) TestFetchMessages() {
	if s.config.CharacterID == "" {
		s.T().Skip("No test character ID provided. Skipping TestFetchMessages.")
	}

	// Create a new chat and send a message
	chat, _, err := s.client.CreateChat(s.config.CharacterID, true)
	s.Require().NoError(err)
	messageText := "Hello, this is a test message."
	_, err = s.client.SendMessage(s.config.CharacterID, chat.ChatID, messageText)
	s.Require().NoError(err)

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)

	// Fetch messages
	turns, _, err := s.client.FetchMessages(chat.ChatID, false, "")
	s.Require().NoError(err, "FetchMessages returned an error")
	s.Assert().NotNil(turns, "Turns should not be nil")
	s.Assert().Greater(len(turns), 0, "There should be at least one turn")
}

func (s *ChatSuite) TestFetchAllMessages() {
	if s.config.CharacterID == "" {
		s.T().Skip("No test character ID provided. Skipping TestFetchAllMessages.")
	}

	// Create a new chat and send messages
	chat, _, err := s.client.CreateChat(s.config.CharacterID, true)
	s.Require().NoError(err)
	messageText := "Hello, this is a test message."
	_, err = s.client.SendMessage(s.config.CharacterID, chat.ChatID, messageText)
	s.Require().NoError(err)

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)

	// Fetch all messages
	turns, err := s.client.FetchAllMessages(chat.ChatID, false)
	s.Require().NoError(err, "FetchAllMessages returned an error")
	s.Assert().NotNil(turns, "Turns should not be nil")
	s.Assert().Greater(len(turns), 0, "There should be at least one turn")
}

// FIXME: Not sure if that is even a functionality? Chat name never part of responses
//func (s *ChatSuite) TestUpdateChatName() {
//	if s.config.CharacterID == "" {
//		s.T().Skip("No test character ID provided. Skipping TestUpdateChatName.")
//	}
//
//	// Create a new chat
//	chat, _, err := s.client.CreateChat(s.config.CharacterID, false)
//	s.Require().NoError(err, "CreateChat returned an error")
//	s.Assert().NotNil(chat, "Chat should not be nil")
//
//	// Update chat name
//	newName := "Updated Chat Name"
//	err = s.client.UpdateChatName(chat.ChatID, newName)
//	s.Require().NoError(err, "UpdateChatName returned an error")
//
//	// Fetch the chat to verify the name
//	updatedChat, err := s.client.FetchChat(chat.ChatID)
//	s.Require().NoError(err, "FetchChat returned an error")
//	s.Assert().Equal(newName, updatedChat.ChatName, "Chat name should be updated")
//}

func (s *ChatSuite) TestArchiveAndUnarchiveChat() {
	if s.config.CharacterID == "" {
		s.T().Skip("No test character ID provided. Skipping TestArchiveAndUnarchiveChat.")
	}

	// Create a new chat
	chat, _, err := s.client.CreateChat(s.config.CharacterID, false)
	s.Require().NoError(err, "CreateChat returned an error")
	s.Assert().NotNil(chat, "Chat should not be nil")

	// Archive chat
	err = s.client.ArchiveChat(chat.ChatID)
	s.Require().NoError(err, "ArchiveChat returned an error")

	// Unarchive chat
	err = s.client.UnarchiveChat(chat.ChatID)
	s.Require().NoError(err, "UnarchiveChat returned an error")
}

func (s *ChatSuite) TestCopyChat() {
	if s.config.CharacterID == "" {
		s.T().Skip("No test character ID provided. Skipping TestCopyChat.")
	}

	// Create a new chat and send a message
	chat, _, err := s.client.CreateChat(s.config.CharacterID, true)
	s.Require().NoError(err)
	messageText := "Hello, this is a test message."
	turn, err := s.client.SendMessage(s.config.CharacterID, chat.ChatID, messageText)
	s.Require().NoError(err)

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)

	// Copy chat up to the last turn ID
	newChatID, err := s.client.CopyChat(chat.ChatID, turn.TurnID)
	s.Require().NoError(err, "CopyChat returned an error")
	s.Assert().NotEmpty(newChatID, "NewChatID should not be empty")

	// Fetch the new chat to verify it exists
	newChat, err := s.client.FetchChat(newChatID)
	s.Require().NoError(err, "FetchChat returned an error")
	s.Assert().Equal(newChatID, newChat.ChatID, "New chat ID should match")
}

func (s *ChatSuite) TestAnotherResponse() {
	if s.config.CharacterID == "" {
		s.T().Skip("No test character ID provided. Skipping TestAnotherResponse.")
	}

	// Create a new chat and send a message
	chat, _, err := s.client.CreateChat(s.config.CharacterID, true)
	s.Require().NoError(err)
	messageText := "Tell me a joke."
	turn, err := s.client.SendMessage(s.config.CharacterID, chat.ChatID, messageText)
	s.Require().NoError(err)

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)

	// Get another response
	newTurn, err := s.client.AnotherResponse(s.config.CharacterID, chat.ChatID, turn.TurnID)
	s.Require().NoError(err, "AnotherResponse returned an error")
	s.Assert().NotNil(newTurn, "New turn should not be nil")

	// Verify that the new response is different
	oldCandidate := turn.Candidates[turn.PrimaryCandidateID]
	newCandidate := newTurn.Candidates[newTurn.PrimaryCandidateID]
	s.Assert().NotEqual(oldCandidate.CandidateID, newCandidate.CandidateID, "Candidate IDs should be different")
}

func (s *ChatSuite) TestEditMessage() {
	if s.config.CharacterID == "" {
		s.T().Skip("No test character ID provided. Skipping TestEditMessage.")
	}

	// Create a new chat and send a message
	chat, _, err := s.client.CreateChat(s.config.CharacterID, true)
	s.Require().NoError(err)
	messageText := "What's the weather today?"
	turn, err := s.client.SendMessage(s.config.CharacterID, chat.ChatID, messageText)
	s.Require().NoError(err)

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)

	// Edit the message
	newText := "What's the weather tomorrow?"
	editedTurn, err := s.client.EditMessage(chat.ChatID, turn.TurnID, turn.PrimaryCandidateID, newText)
	s.Require().NoError(err, "EditMessage returned an error")
	s.Assert().NotNil(editedTurn, "Edited turn should not be nil")

	// Verify that the text has been updated
	editedCandidate := editedTurn.Candidates[editedTurn.PrimaryCandidateID]
	s.Assert().Equal(newText, editedCandidate.Text, "The message text should be updated")
}

func (s *ChatSuite) TestDeleteMessage() {
	if s.config.CharacterID == "" {
		s.T().Skip("No test character ID provided. Skipping TestDeleteMessage.")
	}

	// Create a new chat and send a message
	chat, _, err := s.client.CreateChat(s.config.CharacterID, true)
	s.Require().NoError(err)
	messageText := "This message will be deleted."
	turn, err := s.client.SendMessage(s.config.CharacterID, chat.ChatID, messageText)
	s.Require().NoError(err)

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)

	// Delete the message
	err = s.client.DeleteMessage(chat.ChatID, turn.TurnID)
	s.Require().NoError(err, "DeleteMessage returned an error")

	// Fetch messages and verify the message is deleted
	turns, _, err := s.client.FetchMessages(chat.ChatID, false, "")
	s.Require().NoError(err)
	for _, t := range turns {
		s.Assert().NotEqual(turn.TurnID, t.TurnID, "Deleted turn should not be in the list")
	}
}

func (s *ChatSuite) TestPinAndUnpinMessage() {
	if s.config.CharacterID == "" {
		s.T().Skip("No test character ID provided. Skipping TestPinAndUnpinMessage.")
	}

	// Create a new chat and send a message
	chat, _, err := s.client.CreateChat(s.config.CharacterID, true)
	s.Require().NoError(err)
	messageText := "Please pin this message."
	turn, err := s.client.SendMessage(s.config.CharacterID, chat.ChatID, messageText)
	s.Require().NoError(err)

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)

	// Pin the message
	err = s.client.PinMessage(chat.ChatID, turn.TurnID)
	s.Require().NoError(err, "PinMessage returned an error")

	// Fetch messages and verify the message is pinned
	turns, _, err := s.client.FetchMessages(chat.ChatID, false, "")
	s.Require().NoError(err)
	var pinnedTurn *cai.Turn
	for _, t := range turns {
		if t.TurnID == turn.TurnID {
			pinnedTurn = t
			break
		}
	}
	s.Require().NotNil(pinnedTurn, "Pinned turn should be in the list")
	s.Assert().True(pinnedTurn.IsPinned, "Turn should be pinned")

	// Unpin the message
	err = s.client.UnpinMessage(chat.ChatID, turn.TurnID)
	s.Require().NoError(err, "UnpinMessage returned an error")

	// Fetch messages and verify the message is unpinned
	turns, _, err = s.client.FetchMessages(chat.ChatID, false, "")
	s.Require().NoError(err)
	for _, t := range turns {
		if t.TurnID == turn.TurnID {
			s.Assert().False(t.IsPinned, "Turn should be unpinned")
			break
		}
	}
}

func TestChatSuite(t *testing.T) {
	suite.Run(t, new(ChatSuite))
}
