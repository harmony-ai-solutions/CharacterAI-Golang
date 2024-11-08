package cai

func (s *ClientIntegrationSuite) TestGenerateImage() {
	images, err := s.client.GenerateImage("A futuristic cityscape", 1)
	s.Require().NoError(err, "GenerateImage returned an error")
	s.Assert().NotNil(images, "Images should not be nil")
	s.Assert().Greater(len(images), 0, "Should receive at least one image")

	// Optionally, verify the image URLs
	s.Assert().Contains(images[0], "https://", "Image URL should be valid")
}
