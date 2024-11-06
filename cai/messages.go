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
	Chat Chat `json:"chat"`
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
	Account Account `json:"account"`
}

// FetchMeResponse represents the response from the FetchMe API call
type FetchMeResponse struct {
	User Account `json:"user"`
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
