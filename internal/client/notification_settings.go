package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (c *Client) getNotificationSettings(spaceID string) (*NotificationSettingsResponse, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/spaces/%s/notifications", c.HostURL, spaceID), nil)
	if err != nil {
		return nil, err
	}

	var settings NotificationSettingsResponse

	err = c.doRequestUnmarshal(req, &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil

}

// SpaceDefault returns the space-wide row, the one carrying no room ID.
func (r *NotificationSettingsResponse) SpaceDefault() (*NotificationSettings, error) {

	for _, setting := range r.Settings {
		if setting.RoomID == nil {
			return setting, nil
		}
	}

	return nil, ErrNotFound
}

// Room returns the row overriding the space-wide settings for the given room.
func (r *NotificationSettingsResponse) Room(roomID uuid.UUID) (*NotificationSettings, error) {

	for _, setting := range r.Settings {
		if setting.RoomID != nil && *setting.RoomID == roomID {
			return setting, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) GetSpaceDefaultNotificationSettings(spaceID string) (*NotificationSettings, error) {

	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}

	settings, err := c.getNotificationSettings(spaceID)
	if err != nil {
		return nil, err
	}

	return settings.SpaceDefault()
}

func (c *Client) GetRoomNotificationSettings(spaceID, roomID string) (*NotificationSettings, error) {

	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}

	if roomID == "" {
		return nil, ErrRoomIDRequired
	}

	roomIDParsed, err := uuid.Parse(roomID)

	if err != nil {
		return nil, ErrInvalidRoomID
	}

	settings, err := c.getNotificationSettings(spaceID)
	if err != nil {
		return nil, err
	}

	return settings.Room(roomIDParsed)
}

func (c *Client) UpdateNotificationSettings(spaceID string, settings []NotificationSettingsUpsert) (*NotificationSettingsResponse, error) {

	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}

	if len(settings) == 0 {
		return nil, ErrNotificationSettingsRequired
	}

	reqBody, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v2/spaces/%s/notifications", c.HostURL, spaceID), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	var notificationSettings NotificationSettingsResponse

	err = c.doRequestUnmarshal(req, &notificationSettings)
	if err != nil {
		return nil, err
	}

	return &notificationSettings, nil
}

func (c *Client) DeleteNotificationSettings(spaceID, notificationSettingID string) error {

	if spaceID == "" {
		return ErrSpaceIDRequired
	}

	if notificationSettingID == "" {
		return ErrNotificationSettingIDRequired
	}

	notificationSettingIDParsed, err := uuid.Parse(notificationSettingID)

	if err != nil {
		return ErrInvalidNotificationSettingID
	}

	reqBody, err := json.Marshal([]uuid.UUID{notificationSettingIDParsed})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v2/spaces/%s/notifications", c.HostURL, spaceID), bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil

}
