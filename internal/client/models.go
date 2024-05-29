package client

import "encoding/json"

type SpaceInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoomInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SpaceMember struct {
	Email    string `json:"email"`
	MemberID string `json:"memberID"`
	Role     string `json:"role"`
}

type RoomMember struct {
	SpaceMemberID string `json:"memberID"`
}

type NotificationIntegrations struct {
	Integrations []NotificationIntegration `json:"integrations"`
}

type NotificationChannel struct {
	ID          string                  `json:"id"`
	Enabled     bool                    `json:"enabled"`
	Name        string                  `json:"name"`
	Integration NotificationIntegration `json:"integration"`
	Alarms      string                  `json:"alarms"`
	Rooms       []string                `json:"rooms"`
	Secrets     json.RawMessage         `json:"secrets"`
}

type NotificationIntegration struct {
	ID   string `json:"id"`
	Name string `json:"slug"`
}

type NotificationSlackChannel struct {
	URL string `json:"url"`
}

type NotificationDiscordChannel struct {
	URL           string `json:"url"`
	ChannelParams struct {
		Selection  string `json:"selection"`
		ThreadName string `json:"threadName"`
	} `json:"channelParams"`
}

type NotificationPagerdutyChannel struct {
	AlertEventsURL string `json:"alertEventsURL"`
	IntegrationKey string `json:"integrationKey"`
}

type notificationRequestPayload struct {
	Name          string          `json:"name"`
	IntegrationID string          `json:"integrationID"`
	Alarms        string          `json:"alarms"`
	Rooms         []string        `json:"rooms"`
	Secrets       json.RawMessage `json:"secrets"`
}

type Invitation struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type RoomNodes struct {
	Nodes []RoomNode `json:"nodes"`
}

type RoomNode struct {
	NodeID   string `json:"nd"`
	NodeName string `json:"nm"`
	State    string `json:"state"`
}
