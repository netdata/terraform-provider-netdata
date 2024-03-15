package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetRooms(spaceID string) (*[]RoomInfo, error) {
	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/spaces/%s/rooms", c.HostURL, spaceID), nil)
	if err != nil {
		return nil, err
	}

	var rooms []RoomInfo

	err = c.doRequestUnmarshal(req, &rooms)
	if err != nil {
		return nil, err
	}

	return &rooms, nil
}

func (c *Client) GetRoomByID(id, spaceID string) (*RoomInfo, error) {
	rooms, err := c.GetRooms(spaceID)
	if err != nil {
		return nil, err
	}
	for _, room := range *rooms {
		if room.ID == id {
			return &room, nil
		}
	}
	return nil, ErrNotFound
}

func (c *Client) CreateRoom(spaceID, name, description string) (*RoomInfo, error) {
	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}
	reqBody, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/spaces/%s/rooms", c.HostURL, spaceID), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	var room RoomInfo

	err = c.doRequestUnmarshal(req, &room)
	if err != nil {
		return nil, err
	}

	room.Name = name
	room.Description = description

	return &room, nil
}

func (c *Client) UpdateRoomByID(id, spaceID, name, description string) error {
	if id == "" {
		return fmt.Errorf("id is empty")
	}
	if spaceID == "" {
		return ErrSpaceIDRequired
	}
	reqBody, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/v1/spaces/%s/rooms/%s", c.HostURL, spaceID, id), bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	_, err = c.doRequest(req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteRoomByID(id, spaceID string) error {
	if id == "" {
		return fmt.Errorf("id is empty")
	}
	if spaceID == "" {
		return ErrSpaceIDRequired
	}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v1/spaces/%s/rooms/%s", c.HostURL, spaceID, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
