package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

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
	ID                       string                  `json:"id"`
	Enabled                  bool                    `json:"enabled"`
	Name                     string                  `json:"name"`
	Integration              NotificationIntegration `json:"integration"`
	NotificationOptions      []string                `json:"notification_options"`
	Rooms                    []string                `json:"rooms"`
	Secrets                  json.RawMessage         `json:"secrets"`
	RepeatNotificationMinute int64                   `json:"repeat_notification_min,omitempty"`
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
	Name                     string          `json:"name"`
	IntegrationID            string          `json:"integrationID"`
	NotificationOptions      []string        `json:"notification_options"`
	Rooms                    []string        `json:"rooms"`
	Secrets                  json.RawMessage `json:"secrets"`
	RepeatNotificationMinute int64           `json:"repeat_notification_min,omitempty"`
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
type NodeMembershipRule struct {
	ID          uuid.UUID              `json:"id"`
	SpaceID     uuid.UUID              `json:"spaceID"`
	RoomID      uuid.UUID              `json:"roomID"`
	Clauses     []NodeMembershipClause `json:"clauses"`
	Action      string                 `json:"action"`
	Description string                 `json:"description"`
}
type NodeMembershipClause struct {
	Label    string `json:"label"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
	Negate   bool   `json:"negate"`
}

type SilencingRule struct {
	ID                  uuid.UUID         `json:"id"`
	Name                string            `json:"name"`
	RoomIDs             []uuid.UUID       `json:"room_ids,omitempty"`
	NodeIDs             []uuid.UUID       `json:"node_ids,omitempty"`
	HostLabels          map[string]string `json:"host_labels,omitempty"`
	AlertNames          []string          `json:"alert_names,omitempty"`
	AlertContexts       []string          `json:"alert_contexts,omitempty"`
	AlertInstances      []string          `json:"alert_instances,omitempty"`
	AlertRoles          []string          `json:"alert_roles,omitempty"`
	NotificationOptions []string          `json:"notification_options,omitempty"`
	IntegrationIDs      []uuid.UUID       `json:"integration_ids,omitempty"`
	StartsAt            time.Time         `json:"starts_at"`
	LastsUntil          *time.Time        `json:"lasts_until,omitempty"`
	DeleteOnExpiry      bool              `json:"delete_on_expiry"`
	Disabled            bool              `json:"disabled"`
	RRule               *string           `json:"rrule,omitempty"`
	Timezone            *string           `json:"timezone,omitempty"`
}
