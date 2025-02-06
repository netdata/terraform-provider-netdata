package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) CreateDiscordChannel(spaceID string, commonParams NotificationChannel, discordParams NotificationDiscordChannel) (*NotificationChannel, error) {

	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}

	reqBody := notificationRequestPayload{
		Name:                     commonParams.Name,
		IntegrationID:            commonParams.Integration.ID,
		Rooms:                    commonParams.Rooms,
		NotificationOptions:      commonParams.NotificationOptions,
		RepeatNotificationMinute: commonParams.RepeatNotificationMinute,
	}

	secretsJson, err := json.Marshal(discordParams)
	if err != nil {
		return nil, err
	}

	reqBody.Secrets = json.RawMessage(secretsJson)
	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v2/spaces/%s/channel", c.HostURL, spaceID), bytes.NewReader(jsonReqBody))
	if err != nil {
		return nil, err
	}

	var respNotificationChannel NotificationChannel

	err = c.doRequestUnmarshal(req, &respNotificationChannel)
	if err != nil {
		return nil, err
	}

	err = c.EnableChannelByID(spaceID, respNotificationChannel.ID, commonParams.Enabled)

	if err != nil {
		return nil, err
	}

	respNotificationChannel.Enabled = commonParams.Enabled

	return &respNotificationChannel, nil
}

func (c *Client) UpdateDiscordChannelByID(spaceID string, commonParams NotificationChannel, discordParams NotificationDiscordChannel) (*NotificationChannel, error) {

	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}

	if commonParams.ID == "" {
		return nil, ErrChannelIDRequired
	}

	err := c.EnableChannelByID(spaceID, commonParams.ID, commonParams.Enabled)
	if err != nil {
		return nil, err
	}

	reqBody := notificationRequestPayload{
		Name:                     commonParams.Name,
		Rooms:                    commonParams.Rooms,
		NotificationOptions:      commonParams.NotificationOptions,
		RepeatNotificationMinute: commonParams.RepeatNotificationMinute,
	}

	secretsJson, err := json.Marshal(discordParams)
	if err != nil {
		return nil, err
	}

	reqBody.Secrets = json.RawMessage(secretsJson)
	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/v2/spaces/%s/channel/%s", c.HostURL, spaceID, commonParams.ID), bytes.NewReader(jsonReqBody))
	if err != nil {
		return nil, err
	}

	var respNotificationChannel NotificationChannel

	err = c.doRequestUnmarshal(req, &respNotificationChannel)
	if err != nil {
		return nil, err
	}

	return &respNotificationChannel, nil

}
