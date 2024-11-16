package cai

import (
	"encoding/json"
	"fmt"
	"time"
)

// UserAccount represents a user account.
type UserAccount struct {
	User                     *User    `json:"user"`
	IsHuman                  bool     `json:"is_human"`
	Name                     string   `json:"name"`
	Email                    string   `json:"email"`
	NeedsToAcknowledgePolicy bool     `json:"needs_to_acknowledge_policy"`
	SuspendedUntil           string   `json:"suspended_until"`
	HiddenCharacters         []string `json:"hidden_characters"`
	BlockedUsers             []string `json:"blocked_users"`
	Bio                      string   `json:"bio"`
	Interests                string   `json:"interests"`
	DateOfBirth              string   `json:"date_of_birth"`
}

// User represents the inner "user" object in the JSON.
type User struct {
	Username     string   `json:"username"`
	ID           int64    `json:"id"`
	FirstName    string   `json:"first_name"`
	Account      *Account `json:"account"`
	IsStaff      bool     `json:"is_staff"`
	Subscription bool     `json:"subscription"`
	Entitlements []string `json:"entitlements"`
}

// Account represents the "account" object inside the inner user.
type Account struct {
	Name                     string `json:"name"`
	AvatarType               string `json:"avatar_type"`
	OnboardingComplete       bool   `json:"onboarding_complete"`
	AvatarFileName           string `json:"avatar_file_name"`
	MobileOnboardingComplete bool   `json:"mobile_onboarding_complete,omitempty"`
}

// UnmarshalJSON custom unmarshalling for UserAccount to handle nested account data.
func (a *UserAccount) UnmarshalJSON(data []byte) error {
	type Alias UserAccount // Create an alias to prevent infinite recursion
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}

// PublicUser represents a public user profile.
type PublicUser struct {
	Username         string           `json:"username"`
	Name             string           `json:"name"`
	Bio              string           `json:"bio"`
	AvatarFileName   string           `json:"avatar_file_name"`
	Avatar           *Avatar          `json:"-"`
	NumFollowing     int              `json:"num_following"`
	NumFollowers     int              `json:"num_followers"`
	Characters       []CharacterShort `json:"characters"`
	SubscriptionType string           `json:"subscription_type"`
}

// UnmarshalJSON custom unmarshalling for PublicUser.
func (p *PublicUser) UnmarshalJSON(data []byte) error {
	type Alias PublicUser
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Create Avatar instance if AvatarFileName is provided
	if p.AvatarFileName != "" {
		p.Avatar = &Avatar{FileName: p.AvatarFileName}
	}

	return nil
}

// Persona represents a persona.
type Persona struct {
	PersonaID             string   `json:"external_id"`
	Name                  string   `json:"name"`
	Definition            string   `json:"definition"`
	Greeting              string   `json:"greeting"`
	Description           string   `json:"description"`
	AvatarFileName        string   `json:"avatar_file_name"`
	Avatar                *Avatar  `json:"-"`
	Visibility            string   `json:"visibility"`
	VoiceID               string   `json:"voice_id"`
	AuthorUsername        string   `json:"user__username"`
	Identifier            string   `json:"identifier"`
	Categories            []string `json:"categories"`
	BaseImgPrompt         string   `json:"base_img_prompt"`
	ImgGenEnabled         bool     `json:"img_gen_enabled"`
	Copyable              bool     `json:"copyable"`
	StripImgPromptFromMsg bool     `json:"strip_img_prompt_from_msg"`
	Archived              bool     `json:"archived,omitempty"` // For deleting personas
}

// UnmarshalJSON custom unmarshalling for Persona.
func (p *Persona) UnmarshalJSON(data []byte) error {
	type Alias Persona
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Process Avatar
	if p.AvatarFileName != "" {
		p.Avatar = &Avatar{FileName: p.AvatarFileName}
	}

	return nil
}

// Notification represents a user notification.
type Notification struct {
	NotificationID string    `json:"notification_id"`
	Type           string    `json:"type"`
	Message        string    `json:"message"`
	IsRead         bool      `json:"is_read"`
	CreateTimeStr  string    `json:"create_time"`
	CreateTime     time.Time `json:"-"`
	// Additional fields as needed
}

// UnmarshalJSON custom unmarshalling for Notification.
func (n *Notification) UnmarshalJSON(data []byte) error {
	type Alias Notification
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(n),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse CreateTime
	var err error
	if n.CreateTimeStr != "" {
		n.CreateTime, err = time.Parse(time.RFC3339Nano, n.CreateTimeStr)
		if err != nil {
			return err
		}
	}

	return nil
}

// Avatar represents an avatar image.
type Avatar struct {
	FileName string `json:"file_name"`
}

// GetURL returns the avatar URL.
func (a *Avatar) GetURL(size int, animated bool) string {
	anim := 0
	if animated {
		anim = 1
	}
	return fmt.Sprintf("https://characterai.io/i/%d/static/avatars/%s?webp=true&anim=%d", size, a.FileName, anim)
}

// Voice represents a voice setting.
type Voice struct {
	VoiceID           string       `json:"id"`
	Name              string       `json:"name"`
	Description       string       `json:"description"`
	Gender            string       `json:"gender"`
	Visibility        string       `json:"visibility"`
	PreviewAudioURL   string       `json:"previewAudioURI"`
	PreviewText       string       `json:"previewText"`
	CreatorID         string       `json:"-"`
	CreatorUsername   string       `json:"-"`
	InternalStatus    string       `json:"internalStatus"`
	LastUpdateTime    time.Time    `json:"-"`
	CreatorInfo       *CreatorInfo `json:"creatorInfo,omitempty"`
	LastUpdateTimeStr string       `json:"lastUpdateTime"`
}

// CreatorInfo represents the creator information inside Voice.
type CreatorInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// UnmarshalJSON custom unmarshalling for Voice.
func (v *Voice) UnmarshalJSON(data []byte) error {
	type Alias Voice
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(v),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Set CreatorID and CreatorUsername from CreatorInfo if available
	if v.CreatorInfo != nil {
		v.CreatorID = v.CreatorInfo.ID
		v.CreatorUsername = v.CreatorInfo.Username
	}

	// Parse LastUpdateTime
	var err error
	if v.LastUpdateTimeStr != "" {
		v.LastUpdateTime, err = time.Parse(time.RFC3339Nano, v.LastUpdateTimeStr)
		if err != nil {
			return err
		}
	}

	return nil
}

// CharacterShort represents a summary of a character.
type CharacterShort struct {
	CharacterID     string  `json:"external_id"`
	Title           string  `json:"title"`
	Greeting        string  `json:"greeting"`
	AvatarFileName  string  `json:"avatar_file_name"`
	Avatar          *Avatar `json:"-"`
	Copyable        bool    `json:"copyable"`
	ParticipantName string  `json:"participant__name"`
	AuthorUsername  string  `json:"user__username"`
	NumInteractions int64   `json:"participant__num_interactions"`
	ImgGenEnabled   bool    `json:"img_gen_enabled"`
	Priority        float32 `json:"priority"`
	DefaultVoiceID  *string `json:"default_voice_id"`
	Upvotes         int64   `json:"upvotes"`
	Name            string  `json:"-"`
}

// UnmarshalJSON custom unmarshalling for CharacterShort.
func (cs *CharacterShort) UnmarshalJSON(data []byte) error {
	type Alias CharacterShort
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(cs),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Set Name to ParticipantName
	if cs.Name == "" {
		cs.Name = cs.ParticipantName
	}

	// Create Avatar instance if AvatarFileName is provided
	if cs.AvatarFileName != "" {
		cs.Avatar = &Avatar{FileName: cs.AvatarFileName}
	}

	return nil
}

// Character represents a CharacterAI character.
type Character struct {
	ExternalID      string                 `json:"external_id"`
	Title           string                 `json:"title"`
	Name            string                 `json:"name"`
	ParticipantName string                 `json:"participant__name"`
	Visibility      string                 `json:"visibility"`
	Greeting        string                 `json:"greeting"`
	Description     string                 `json:"description"`
	Definition      string                 `json:"definition"`
	Upvotes         int64                  `json:"upvotes"`
	AuthorUsername  string                 `json:"user__username"`
	NumInteractions int64                  `json:"participant__num_interactions"`
	AvatarFileName  string                 `json:"avatar_file_name"`
	Avatar          *Avatar                `json:"-"`
	Copyable        bool                   `json:"copyable"`
	Identifier      string                 `json:"identifier"`
	ImgGenEnabled   bool                   `json:"img_gen_enabled"`
	BaseImgPrompt   string                 `json:"base_img_prompt"`
	ImgPromptRegex  string                 `json:"img_prompt_regex"`
	StripImgPrompt  bool                   `json:"strip_img_prompt_from_msg"`
	StarterPrompts  map[string]interface{} `json:"starter_prompts"`
	CommentsEnabled *bool                  `json:"comments_enabled"`
	CreatorName     string                 `json:"participant__user__username"`
	VoiceID         string                 `json:"voice_id"`
	DefaultVoiceID  string                 `json:"default_voice_id"`
	Songs           []string               `json:"songs"`
}

// UnmarshalJSON custom unmarshalling for Character.
func (c *Character) UnmarshalJSON(data []byte) error {
	type Alias Character
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Choose Name from ParticipantName or Name
	if c.Name == "" {
		c.Name = c.ParticipantName
	}

	// Create Avatar instance if AvatarFileName is provided
	if c.AvatarFileName != "" {
		c.Avatar = &Avatar{FileName: c.AvatarFileName}
	}

	return nil
}

type CharacterSearchResult struct {
	DocumentID              string  `json:"document_id"`
	ExternalID              string  `json:"external_id"`
	Title                   string  `json:"title"`
	Greeting                string  `json:"greeting"`
	AvatarFileName          string  `json:"avatar_file_name"`
	Avatar                  *Avatar `json:"-"`
	Visibility              string  `json:"visibility"`
	ParticipantName         string  `json:"participant__name"`
	ParticipantInteractions float64 `json:"participant__num_interactions"`
	AuthorUsername          string  `json:"user__username"`
	Priority                float64 `json:"priority"`
	SearchScore             float64 `json:"search_score"`
	Name                    string  `json:"-"`
}

func (csr *CharacterSearchResult) UnmarshalJSON(data []byte) error {
	type Alias CharacterSearchResult
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(csr),
	}

	// Unmarshal into the auxiliary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Set Name to ParticipantName if Name is empty
	if csr.Name == "" {
		csr.Name = csr.ParticipantName
	}

	// Create Avatar instance if AvatarFileName is provided
	if csr.AvatarFileName != "" {
		csr.Avatar = &Avatar{FileName: csr.AvatarFileName}
	}

	return nil
}

type TurnKey struct {
	ChatID string `json:"chat_id"`
	TurnID string `json:"turn_id"`
}

// Turn represents a chat turn.
type Turn struct {
	TurnKey            TurnKey                   `json:"turn_key"`
	ChatID             string                    `json:"-"`
	TurnID             string                    `json:"-"`
	CreateTimeStr      string                    `json:"create_time"`
	CreateTime         time.Time                 `json:"-"`
	LastUpdateTimeStr  string                    `json:"last_update_time"`
	LastUpdateTime     time.Time                 `json:"-"`
	State              string                    `json:"state"`
	Author             AuthorInfo                `json:"author"`
	CandidatesList     []TurnCandidate           `json:"candidates"`
	Candidates         map[string]*TurnCandidate `json:"-"`
	PrimaryCandidateID string                    `json:"primary_candidate_id"`
	IsPinned           bool                      `json:"is_pinned"`
}

// UnmarshalJSON custom unmarshalling for Turn.
func (t *Turn) UnmarshalJSON(data []byte) error {
	type Alias Turn
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse times
	var err error
	if t.CreateTimeStr != "" {
		t.CreateTime, err = time.Parse(time.RFC3339Nano, t.CreateTimeStr)
		if err != nil {
			return err
		}
	}
	if t.LastUpdateTimeStr != "" {
		t.LastUpdateTime, err = time.Parse(time.RFC3339Nano, t.LastUpdateTimeStr)
		if err != nil {
			return err
		}
	}

	// Build Candidates map
	t.Candidates = make(map[string]*TurnCandidate)
	for i := range t.CandidatesList {
		c := &t.CandidatesList[i]
		t.Candidates[c.CandidateID] = c
	}

	// Set ChatID and TurnID from TurnKey
	t.ChatID = t.TurnKey.ChatID
	t.TurnID = t.TurnKey.TurnID

	return nil
}

// AuthorInfo represents the author of a turn
type AuthorInfo struct {
	AuthorID string `json:"author_id"`
	Name     string `json:"name"`
	IsHuman  bool   `json:"is_human"`
}

// TurnCandidate represents a candidate response.
type TurnCandidate struct {
	CandidateID   string    `json:"candidate_id"`
	Text          string    `json:"raw_content"`
	IsFinal       bool      `json:"is_final"`
	IsFiltered    bool      `json:"safety_truncated"`
	CreateTimeStr string    `json:"create_time"`
	CreateTime    time.Time `json:"-"`
}

// UnmarshalJSON custom unmarshalling for TurnCandidate.
func (tc *TurnCandidate) UnmarshalJSON(data []byte) error {
	type Alias TurnCandidate
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(tc),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse CreateTime
	var err error
	if tc.CreateTimeStr != "" {
		tc.CreateTime, err = time.Parse(time.RFC3339Nano, tc.CreateTimeStr)
		if err != nil {
			return err
		}
	}

	return nil
}

// Chat represents a chat session.
type Chat struct {
	ChatID             string    `json:"chat_id"`
	CharacterID        string    `json:"character_id"`
	CreatorID          string    `json:"creator_id"`
	CreateTimeStr      string    `json:"create_time"`
	CreateTime         time.Time `json:"-"`
	State              string    `json:"state"`
	ChatType           string    `json:"type"`
	Visibility         string    `json:"visibility"`
	PreviewTurns       []*Turn   `json:"preview_turns,omitempty"`
	ChatName           string    `json:"name,omitempty"`
	CharacterName      string    `json:"character_name,omitempty"`
	CharacterAvatarURI string    `json:"character_avatar_uri,omitempty"`
	CharacterAvatar    *Avatar   `json:"-"`
}

// UnmarshalJSON custom unmarshalling for Chat.
func (c *Chat) UnmarshalJSON(data []byte) error {
	type Alias Chat
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse CreateTime
	var err error
	if c.CreateTimeStr != "" {
		c.CreateTime, err = time.Parse(time.RFC3339Nano, c.CreateTimeStr)
		if err != nil {
			return err
		}
	}

	// Lowercase the visibility
	c.Visibility = c.Visibility

	// Create CharacterAvatar instance if CharacterAvatarURI is provided
	if c.CharacterAvatarURI != "" {
		c.CharacterAvatar = &Avatar{FileName: c.CharacterAvatarURI}
	}

	return nil
}

// ChatHistory represents the chat history.
type ChatHistory struct {
	ChatID             string           `json:"external_id"`
	CreateTimeStr      string           `json:"created"`
	CreateTime         time.Time        `json:"-"`
	LastInteractionStr string           `json:"last_interaction"`
	LastInteraction    time.Time        `json:"-"`
	PreviewMessages    []HistoryMessage `json:"msgs"`
}

// UnmarshalJSON custom unmarshalling for ChatHistory.
func (ch *ChatHistory) UnmarshalJSON(data []byte) error {
	type Alias ChatHistory
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(ch),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse times
	var err error
	if ch.CreateTimeStr != "" {
		ch.CreateTime, err = time.Parse(time.RFC3339Nano, ch.CreateTimeStr)
		if err != nil {
			return err
		}
	}
	if ch.LastInteractionStr != "" {
		ch.LastInteraction, err = time.Parse(time.RFC3339Nano, ch.LastInteractionStr)
		if err != nil {
			return err
		}
	}

	return nil
}

type HistoryMessage struct {
	UUID              string `json:"uuid"`
	ID                string `json:"id"`
	Text              string `json:"text"`
	Src               string `json:"src"`
	Tgt               string `json:"tgt"`
	IsAlternative     bool   `json:"is_alternative"`
	ImageRelativePath string `json:"image_rel_path"`
}

// PreviousAnnotations represents the annotations used in the request.
type PreviousAnnotations map[string]int

// Function to generate previous annotations with zeros.
func generatePreviousAnnotations() PreviousAnnotations {
	return PreviousAnnotations{
		"bad_memory":           0,
		"boring":               0,
		"ends_chat_early":      0,
		"funny":                0,
		"helpful":              0,
		"inaccurate":           0,
		"interesting":          0,
		"long":                 0,
		"not_bad_memory":       0,
		"not_boring":           0,
		"not_ends_chat_early":  0,
		"not_funny":            0,
		"not_helpful":          0,
		"not_inaccurate":       0,
		"not_interesting":      0,
		"not_long":             0,
		"not_out_of_character": 0,
		"not_repetitive":       0,
		"not_short":            0,
		"out_of_character":     0,
		"repetitive":           0,
		"short":                0,
	}
}

// Settings represents the user's settings.
type Settings struct {
	DefaultPersonaID string            `json:"default_persona_id"`
	EnableTTS        bool              `json:"enable_tts"`
	PersonaOverrides map[string]string `json:"personaOverrides"`
}

// UnmarshalJSON custom unmarshal to capture additional fields.
func (s *Settings) UnmarshalJSON(data []byte) error {
	type Alias Settings
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	// Unmarshal known fields
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}
