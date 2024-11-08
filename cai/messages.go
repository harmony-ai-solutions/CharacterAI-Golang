package cai

import "encoding/json"

type WebSocketMessage struct {
	Command   string      `json:"command"`
	OriginID  string      `json:"origin_id,omitempty"`
	RequestID string      `json:"request_id"`
	Payload   interface{} `json:"payload"`
}

// CreateAndGenerateTurnPayload represents the payload for creating and generating a turn.
type CreateAndGenerateTurnPayload struct {
	CharacterID         string              `json:"character_id"`
	NumCandidates       int                 `json:"num_candidates"`
	PreviousAnnotations PreviousAnnotations `json:"previous_annotations"`
	SelectedLanguage    string              `json:"selected_language"`
	TTSEnabled          bool                `json:"tts_enabled"`
	Turn                TurnPayload         `json:"turn"`
	UserName            string              `json:"user_name"`
}

type TurnPayload struct {
	Author             AuthorPayload      `json:"author"`
	Candidates         []CandidatePayload `json:"candidates"`
	PrimaryCandidateID string             `json:"primary_candidate_id"`
	TurnKey            TurnKey            `json:"turn_key"`
}

type AuthorPayload struct {
	AuthorID string `json:"author_id"`
	IsHuman  bool   `json:"is_human"`
	Name     string `json:"name,omitempty"`
}

type CandidatePayload struct {
	CandidateID string `json:"candidate_id"`
	RawContent  string `json:"raw_content"`
}

type WebSocketResponse struct {
	Command string          `json:"command"`
	Payload json.RawMessage `json:"payload,omitempty"`
	Comment string          `json:"comment,omitempty"`
}

type TurnResponsePayload struct {
	Turn Turn `json:"turn"`
}

type CreateChatPayload struct {
	Chat         ChatPayload `json:"chat"`
	WithGreeting bool        `json:"with_greeting"`
}

type ChatPayload struct {
	ChatID      string `json:"chat_id"`
	CreatorID   string `json:"creator_id"`
	Visibility  string `json:"visibility"`
	CharacterID string `json:"character_id"`
	Type        string `json:"type"`
}

type CreateChatResponsePayload struct {
	Chat      Chat   `json:"chat"`
	Command   string `json:"command"`
	RequestID string `json:"request_id"`
}

type UpdatePrimaryCandidatePayload struct {
	CandidateID string  `json:"candidate_id"`
	TurnKey     TurnKey `json:"turn_key"`
}

type EditTurnCandidatePayload struct {
	TurnKey                TurnKey `json:"turn_key"`
	CurrentCandidateID     string  `json:"current_candidate_id"`
	NewCandidateRawContent string  `json:"new_candidate_raw_content"`
}

type RemoveTurnsPayload struct {
	ChatID  string   `json:"chat_id"`
	TurnIDs []string `json:"turn_ids"`
}

type SetTurnPinPayload struct {
	IsPinned bool    `json:"is_pinned"`
	TurnKey  TurnKey `json:"turn_key"`
}

type FetchChatsResponse struct {
	Chats []*Chat `json:"chats"`
}

type FetchHistoriesRequest struct {
	ExternalID string `json:"external_id"`
	Number     int    `json:"number"`
}

type FetchHistoriesResponse struct {
	Histories []ChatHistory `json:"histories"`
}

type CopyChatRequest struct {
	EndTurnID string `json:"end_turn_id"`
}

type CopyChatResponse struct {
	NewChatID string `json:"new_chat_id"`
}

// Meta information for pagination
type Meta struct {
	NextToken string `json:"next_token"`
}

type FetchMessagesResponse struct {
	Meta  Meta   `json:"meta"`
	Turns []Turn `json:"turns"`
}

type MarkNotificationsReadPayload struct {
	NotificationIDs []string `json:"notification_ids"`
}

// UpdateProfilePayload represents the payload for updating the user's profile.
type UpdateProfilePayload struct {
	AvatarType    string `json:"avatar_type"`
	Bio           string `json:"bio"`
	Name          string `json:"name"`
	Username      string `json:"username"`
	AvatarRelPath string `json:"avatar_rel_path,omitempty"`
	// Additional fields that might be required
	Email            string `json:"email,omitempty"`
	FirstName        string `json:"first_name,omitempty"`
	IsHuman          bool   `json:"is_human,omitempty"`
	SubscriptionType string `json:"subscription_type,omitempty"`
}

type NotificationsResponsePayload struct {
	Notifications []Notification `json:"notifications"`
}

type AccountInfoResponsePayload struct {
	Account UserAccount `json:"account"`
}

// FetchMeResponse represents the response from the FetchMe API call
type FetchMeResponse struct {
	User UserAccount `json:"user"`
}

type GenerateTurnCandidatePayload struct {
	CharacterID         string              `json:"character_id"`
	TTSEnabled          bool                `json:"tts_enabled"`
	PreviousAnnotations PreviousAnnotations `json:"previous_annotations"`
	SelectedLanguage    string              `json:"selected_language"`
	UserName            string              `json:"user_name"`
	TurnKey             TurnKey             `json:"turn_key"`
}

// CreatePersonaPayload represents the payload for creating a persona.
type CreatePersonaPayload struct {
	Name                  string   `json:"name"`
	Title                 string   `json:"title"`
	Definition            string   `json:"definition"`
	Greeting              string   `json:"greeting"`
	Description           string   `json:"description"`
	Visibility            string   `json:"visibility"`
	AvatarFileName        string   `json:"avatar_file_name"`
	AvatarRelPath         string   `json:"avatar_rel_path,omitempty"`
	VoiceID               string   `json:"voice_id"`
	Identifier            string   `json:"identifier"`
	Categories            []string `json:"categories"`
	BaseImgPrompt         string   `json:"base_img_prompt"`
	ImgGenEnabled         bool     `json:"img_gen_enabled"`
	Copyable              bool     `json:"copyable"`
	StripImgPromptFromMsg bool     `json:"strip_img_prompt_from_msg"`
}

// EditPersonaPayload represents the payload for editing a persona.
type EditPersonaPayload struct {
	ExternalID            string   `json:"external_id"`
	Name                  string   `json:"name"`
	Title                 string   `json:"title"`
	Definition            string   `json:"definition"`
	Greeting              string   `json:"greeting"`
	Description           string   `json:"description"`
	Visibility            string   `json:"visibility"`
	AvatarFileName        string   `json:"avatar_file_name"`
	AvatarRelPath         string   `json:"avatar_rel_path,omitempty"`
	VoiceID               string   `json:"voice_id"`
	Identifier            string   `json:"identifier"`
	Categories            []string `json:"categories"`
	BaseImgPrompt         string   `json:"base_img_prompt"`
	ImgGenEnabled         bool     `json:"img_gen_enabled"`
	Copyable              bool     `json:"copyable"`
	StripImgPromptFromMsg bool     `json:"strip_img_prompt_from_msg"`
	Archived              bool     `json:"archived,omitempty"` // For deleting personas
}

type CreatePersonaResult struct {
	Status  string   `json:"status"`
	Persona *Persona `json:"persona"`
	Error   string   `json:"error"`
}

type EditPersonaResult struct {
	Status  string   `json:"status"`
	Persona *Persona `json:"persona"`
	Error   string   `json:"error"`
}

type DeletePersonaResult struct {
	Status  string   `json:"status"`
	Persona *Persona `json:"persona"`
	Error   string   `json:"error"`
}

// CreateCharacterPayload represents the payload for creating a character.
type CreateCharacterPayload struct {
	Name                  string   `json:"name"`
	Title                 string   `json:"title"`
	Definition            string   `json:"definition"`
	Greeting              string   `json:"greeting"`
	Description           string   `json:"description"`
	Visibility            string   `json:"visibility"`
	AvatarRelPath         string   `json:"avatar_rel_path,omitempty"`
	BaseImgPrompt         string   `json:"base_img_prompt"`
	Categories            []string `json:"categories"`
	Copyable              bool     `json:"copyable"`
	DefaultVoiceID        string   `json:"default_voice_id"`
	Identifier            string   `json:"identifier"`
	ImgGenEnabled         bool     `json:"img_gen_enabled"`
	StripImgPromptFromMsg bool     `json:"strip_img_prompt_from_msg"`
	VoiceID               string   `json:"voice_id"`
}

// CreateCharacterResult represents the result of creating a character.
type CreateCharacterResult struct {
	Status    string     `json:"status"`
	Character *Character `json:"character"`
	Error     string     `json:"error"`
}

// EditCharacterPayload represents the payload for editing a character.
type EditCharacterPayload struct {
	ExternalID            string   `json:"external_id"`
	Name                  string   `json:"name"`
	Title                 string   `json:"title"`
	Definition            string   `json:"definition"`
	Greeting              string   `json:"greeting"`
	Description           string   `json:"description"`
	Visibility            string   `json:"visibility"`
	AvatarRelPath         string   `json:"avatar_rel_path,omitempty"`
	BaseImgPrompt         string   `json:"base_img_prompt"`
	Categories            []string `json:"categories"`
	Copyable              bool     `json:"copyable"`
	DefaultVoiceID        string   `json:"default_voice_id"`
	ImgGenEnabled         bool     `json:"img_gen_enabled"`
	StripImgPromptFromMsg bool     `json:"strip_img_prompt_from_msg"`
	VoiceID               string   `json:"voice_id"`
	Archived              bool     `json:"archived"`
}

// EditCharacterResult represents the result of editing a character.
type EditCharacterResult struct {
	Status    string     `json:"status"`
	Character *Character `json:"character"`
	Error     string     `json:"error"`
}

// CharacterInfoPayload represents the request payload for fetching character info.
type CharacterInfoPayload struct {
	ExternalID string `json:"external_id"`
}

// CharacterInfoResponse represents the response payload for fetching character info.
type CharacterInfoResponse struct {
	Status    string     `json:"status"`
	Character *Character `json:"character"`
	Error     string     `json:"error"`
}

// CharacterVotePayload represents the request payload for voting on a character.
type CharacterVotePayload struct {
	ExternalID string `json:"external_id"`
	Vote       *bool  `json:"vote"` // Pointer to bool to handle nil (removing vote)
}

// CharacterVoteResponse represents the response payload after voting.
type CharacterVoteResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// SearchCreatorsResponse represents the response payload for searching creators.
type SearchCreatorsResponse struct {
	Status   string    `json:"status"`
	Creators []Creator `json:"creators"`
	HasMore  bool      `json:"has_more"`
	Error    string    `json:"error"`
}

// Creator represents a creator in the system.
type Creator struct {
	Name string `json:"name"`
}

// SearchVoicesResponse represents the response for searching voices.
type SearchVoicesResponse struct {
	Voices []*Voice `json:"voices"`
}

// FetchVoiceResponse represents the response for fetching a voice.
type FetchVoiceResponse struct {
	Voice *Voice `json:"voice"`
}

// UploadVoiceResponse represents the response after uploading or editing a voice.
type UploadVoiceResponse struct {
	Voice *Voice `json:"voice"`
}

// VoiceUpdatePayload represents the payload for updating a voice.
type VoiceUpdatePayload struct {
	Voice VoiceInfo `json:"voice"`
}

// VoiceInfo represents the voice information for updates.
type VoiceInfo struct {
	AudioSourceType string       `json:"audioSourceType"`
	BackendID       string       `json:"backendId"`
	BackendProvider string       `json:"backendProvider"`
	CreatorInfo     *CreatorInfo `json:"creatorInfo"`
	Description     string       `json:"description"`
	Gender          string       `json:"gender"`
	ID              string       `json:"id"`
	InternalStatus  string       `json:"internalStatus"`
	LastUpdateTime  string       `json:"lastUpdateTime"`
	Name            string       `json:"name"`
	PreviewAudioURI string       `json:"previewAudioURI"`
	PreviewText     string       `json:"previewText"`
	Visibility      string       `json:"visibility"`
}

// GenerateSpeechPayload represents the payload for generating speech.
type GenerateSpeechPayload struct {
	CandidateID string `json:"candidateId"`
	RoomID      string `json:"roomId"`
	TurnID      string `json:"turnId"`
	VoiceID     string `json:"voiceId"`
}

// GenerateSpeechResponse represents the response from generating speech.
type GenerateSpeechResponse struct {
	ReplayURL string `json:"replayUrl"`
}

// ErrorResponse represents an error response from the API.
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// FetchUserRequest represents the request payload for fetching a user's public information.
type FetchUserRequest struct {
	Username string `json:"username"`
}

// FetchUserResponse represents the response payload when fetching a user's public information.
type FetchUserResponse struct {
	PublicUser *PublicUser `json:"public_user"`
}

// FollowUserRequest represents the request payload for following/unfollowing a user.
type FollowUserRequest struct {
	Username string `json:"username"`
}

// FollowUserResponse represents the response payload after following/unfollowing a user.
type FollowUserResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// GenerateImageRequest represents the request payload for generating images.
type GenerateImageRequest struct {
	Prompt        string `json:"prompt"`
	NumCandidates int    `json:"num_candidates"`
	ModelVersion  string `json:"model_version"`
}

type GenerateImageResponse struct {
	Result []ImageResult `json:"result"`
}

// ImageResult represents a single image result.
type ImageResult struct {
	URL string `json:"url"`
}

// UploadAvatarRequest represents the request payload when uploading an avatar.
type UploadAvatarRequest map[string]AvatarRequestItem

// AvatarRequestItem represents a single request item in the avatar upload payload.
type AvatarRequestItem struct {
	JSON AvatarJSON `json:"json"`
}

// AvatarJSON contains the image data URL.
type AvatarJSON struct {
	ImageDataURL string `json:"imageDataUrl"`
}

// UploadAvatarResponse represents the response when uploading an avatar.
type UploadAvatarResponse struct {
	Result struct {
		Data struct {
			JSON string `json:"json"`
		} `json:"data"`
	} `json:"result"`
}

// UploadVoiceMetadata represents the metadata for uploading a voice.
type UploadVoiceMetadata struct {
	Voice UploadVoiceInfo `json:"voice"`
}

// UploadVoiceInfo represents the voice information for upload.
type UploadVoiceInfo struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Gender          string `json:"gender"`
	Visibility      string `json:"visibility"`
	PreviewText     string `json:"previewText"`
	AudioSourceType string `json:"audioSourceType"`
}

// UpdateChatNamePayload represents the payload for updating the chat name.
type UpdateChatNamePayload struct {
	Name string `json:"name"`
}

// SetVoicePayload represents the payload for setting a voice override.
type SetVoicePayload struct {
	VoiceID string `json:"voice_id"`
}

// SetVoiceResponse represents the response after setting a voice override.
type SetVoiceResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
