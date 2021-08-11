package slackutils

import "encoding/json"

type Type string

const (
	TypeEventCallBack   = "event_callback"
	TypeURLVerification = "url_verification"
)

// GeneralRequest uses as general case that we don't know what type of request
type GeneralRequest struct {
	Token string `json:"token"`
	Type  Type   `json:"type"`
}

type URLVerificationRequest struct {
	GeneralRequest
	Challenge string `json:"challenge"`
}

type MessageEventRequest struct {
	GeneralRequest
	TeamID   string `json:"team_id"`
	APIAppID string `json:"api_app_id"`
	Event    struct {
		BotID       string `json:"bot_id"`
		Type        string `json:"type"`
		Channel     string `json:"channel"`
		User        string `json:"user"`
		Text        string `json:"text"`
		Ts          string `json:"ts"`
		EventTs     string `json:"event_ts"`
		ChannelType string `json:"channel_type"`
	} `json:"event"`
	AuthedTeams []string `json:"authed_teams"`
	EventID     string   `json:"event_id"`
	EventTime   int      `json:"event_time"`
	ChannelID   string   `json:"channel_id"`
	ChannelName string   `json:"channel_name"`
}

func ParseGeneralRequest(raw []byte) (result *GeneralRequest, err error) {
	var request GeneralRequest
	return &request, json.Unmarshal(raw, &request)
}

func ParseURLVerificationRequest(raw []byte) (result *URLVerificationRequest, err error) {
	var request URLVerificationRequest
	return &request, json.Unmarshal(raw, &request)
}

func ParseMessageEventRequest(raw []byte) (result *MessageEventRequest, err error) {
	var request MessageEventRequest
	return &request, json.Unmarshal(raw, &request)
}

func (r GeneralRequest) ValidateToken(validToken string) bool {
	return r.Token == validToken
}

func (r MessageEventRequest) IsMessageFromBot() bool {
	return len(r.Event.BotID) > 0
}
