package cai

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type VoiceSuite struct {
	BaseSuite
}

func (s *VoiceSuite) TestFetchVoice() {
	// Replace with your actual test code
	voiceID := "test-voice-id"
	voice, err := s.client.FetchVoice(voiceID)
	s.Require().NoError(err, "FetchVoice returned an error")
	s.Assert().NotNil(voice, "Voice should not be nil")
	s.Assert().Equal(voiceID, voice.VoiceID, "VoiceID should match")

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)
}

func TestVoiceSuite(t *testing.T) {
	suite.Run(t, new(VoiceSuite))
}
