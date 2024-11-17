package cai

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ImageSuite struct {
	BaseSuite
}

func (s *ImageSuite) TestGenerateImage() {
	images, err := s.client.GenerateImage("A futuristic cityscape", 1)
	s.Require().NoError(err, "GenerateImage returned an error")
	s.Assert().NotNil(images, "Images should not be nil")
	s.Assert().Greater(len(images), 0, "Should receive at least one image")

	// Optionally, verify the image URLs
	s.Assert().Contains(images[0], "https://", "Image URL should be valid")

	// Pause to avoid rate limiting
	time.Sleep(1 * time.Second)
}

func TestImageSuite(t *testing.T) {
	suite.Run(t, new(ImageSuite))
}
