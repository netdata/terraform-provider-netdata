package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetNotificationChannelByIDAndType(spaceID, channelID, typeName string) (*NotificationChannel, error) {

	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}

	if channelID == "" {
		return nil, ErrChannelIDRequired
	}

	channels, err := c.GetNotificationChannelByType(spaceID, typeName)
	if err != nil {
		return nil, err
	}

	var found bool
	for _, channel := range *channels {
		if channel.ID != channelID {
			continue
		}
		found = true
		break
	}

	if !found {
		return nil, ErrNotFound
	}

	// if the channel is found, get the detailed channel information
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/spaces/%s/channel/%s", c.HostURL, spaceID, channelID), nil)
	if err != nil {
		return nil, err
	}

	var channelDetailed NotificationChannel

	err = c.doRequestUnmarshal(req, &channelDetailed)
	if err != nil {
		return nil, err
	}

	return &channelDetailed, nil

}

func (c *Client) GetNotificationIntegrationByType(spaceID, typeName string) (*NotificationIntegration, error) {

	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/spaces/%s/integrations", c.HostURL, spaceID), nil)
	if err != nil {
		return nil, err
	}

	var integrations NotificationIntegrations

	err = c.doRequestUnmarshal(req, &integrations)
	if err != nil {
		return nil, err
	}

	for _, integration := range integrations.Integrations {
		if strings.EqualFold(integration.Name, typeName) {
			integration.Name = strings.ToLower(integration.Name)
			return &integration, nil
		}
	}

	return nil, ErrNotFound

}

func (c *Client) GetNotificationChannelByType(spaceID, typeName string) (*[]NotificationChannel, error) {

	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/spaces/%s/channel", c.HostURL, spaceID), nil)
	if err != nil {
		return nil, err
	}

	var channels []NotificationChannel

	err = c.doRequestUnmarshal(req, &channels)
	if err != nil {
		return nil, err
	}

	var result []NotificationChannel
	for _, channel := range channels {
		if strings.EqualFold(channel.Integration.Name, typeName) {
			channel.Integration.Name = strings.ToLower(channel.Integration.Name)
			result = append(result, channel)
		}
	}

	if len(result) > 0 {
		return &result, nil
	}

	return nil, ErrNotFound

}

func (c *Client) EnableChannelByID(spaceID, channelID string, enabled bool) error {

	if spaceID == "" {
		return ErrSpaceIDRequired
	}

	if channelID == "" {
		return ErrChannelIDRequired
	}

	reqBody, err := json.Marshal(map[string]bool{
		"enabled": enabled,
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/v2/spaces/%s/channel/%s", c.HostURL, spaceID, channelID), bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteChannelByID(spaceID, channelID string) error {

	if spaceID == "" {
		return ErrSpaceIDRequired
	}

	if channelID == "" {
		return ErrChannelIDRequired
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v2/spaces/%s/channel/%s", c.HostURL, spaceID, channelID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
